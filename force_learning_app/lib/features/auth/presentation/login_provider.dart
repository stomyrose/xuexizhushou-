import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../data/auth_repository.dart';

class LoginState {
  final bool isLoading;
  final String? error;

  LoginState({this.isLoading = false, this.error});

  LoginState copyWith({bool? isLoading, String? error}) {
    return LoginState(
      isLoading: isLoading ?? this.isLoading,
      error: error,
    );
  }
}

class LoginNotifier extends StateNotifier<LoginState> {
  LoginNotifier() : super(LoginState());

  Future<bool> login({String? email, String? phone, required String password}) async {
    state = state.copyWith(isLoading: true, error: null);
    
    try {
      await authRepository.login(
        email: email,
        phone: phone,
        password: password,
      );
      state = state.copyWith(isLoading: false);
      return true;
    } catch (e) {
      state = state.copyWith(isLoading: false, error: e.toString());
      return false;
    }
  }

  void clearError() {
    state = state.copyWith(error: null);
  }
}

final loginProvider = StateNotifierProvider<LoginNotifier, LoginState>((ref) {
  return LoginNotifier();
});
