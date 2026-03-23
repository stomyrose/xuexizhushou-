import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../data/subscription_repository.dart';
import '../../../core/api/models.dart';

class SubscriptionState {
  final List<SubscriptionPlan> plans;
  final SubscriptionPlan? currentSubscription;
  final bool isLoading;
  final String? error;
  final bool purchaseSuccess;

  SubscriptionState({
    this.plans = const [],
    this.currentSubscription,
    this.isLoading = false,
    this.error,
    this.purchaseSuccess = false,
  });

  SubscriptionState copyWith({
    List<SubscriptionPlan>? plans,
    SubscriptionPlan? currentSubscription,
    bool? isLoading,
    String? error,
    bool? purchaseSuccess,
  }) {
    return SubscriptionState(
      plans: plans ?? this.plans,
      currentSubscription: currentSubscription ?? this.currentSubscription,
      isLoading: isLoading ?? this.isLoading,
      error: error,
      purchaseSuccess: purchaseSuccess ?? this.purchaseSuccess,
    );
  }
}

class SubscriptionNotifier extends StateNotifier<SubscriptionState> {
  SubscriptionNotifier() : super(SubscriptionState());

  Future<void> loadPlans() async {
    state = state.copyWith(isLoading: true, error: null);
    try {
      final plans = await subscriptionRepository.getPlans();
      final current = await subscriptionRepository.getCurrentSubscription();
      state = state.copyWith(
        plans: plans,
        currentSubscription: current,
        isLoading: false,
      );
    } catch (e) {
      state = state.copyWith(isLoading: false, error: e.toString());
    }
  }

  Future<bool> purchase(String planId) async {
    state = state.copyWith(isLoading: true, error: null, purchaseSuccess: false);
    try {
      await subscriptionRepository.purchase(planId);
      await loadPlans();
      state = state.copyWith(isLoading: false, purchaseSuccess: true);
      return true;
    } catch (e) {
      state = state.copyWith(isLoading: false, error: e.toString());
      return false;
    }
  }
}

final subscriptionProvider = StateNotifierProvider<SubscriptionNotifier, SubscriptionState>((ref) {
  return SubscriptionNotifier();
});
