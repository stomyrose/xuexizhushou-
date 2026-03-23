import 'package:flutter/foundation.dart';
import '../../../core/api/api_client.dart';
import '../../../core/api/models.dart';

class SubscriptionRepository {
  Future<List<SubscriptionPlan>> getPlans() async {
    try {
      final response = await apiClient.get('/subscriptions/plans');
      final apiResponse = ApiResponse.fromJson(response.data);
      
      if (!apiResponse.isSuccess) {
        throw Exception(apiResponse.message);
      }

      final List<dynamic> data = apiResponse.data ?? [];
      return data.map((e) => SubscriptionPlan.fromJson(e)).toList();
    } catch (e) {
      debugPrint('Get plans error: $e');
      rethrow;
    }
  }

  Future<SubscriptionPlan?> getCurrentSubscription() async {
    try {
      final response = await apiClient.get('/subscriptions/current');
      final apiResponse = ApiResponse.fromJson(response.data);
      
      if (!apiResponse.isSuccess) {
        return null;
      }

      return SubscriptionPlan.fromJson(apiResponse.data);
    } catch (e) {
      debugPrint('Get current subscription error: $e');
      return null;
    }
  }

  Future<bool> purchase(String planId) async {
    try {
      final response = await apiClient.post('/subscriptions/purchase', data: {
        'plan_id': planId,
      });
      
      final apiResponse = ApiResponse.fromJson(response.data);
      return apiResponse.isSuccess;
    } catch (e) {
      debugPrint('Purchase error: $e');
      rethrow;
    }
  }
}

final subscriptionRepository = SubscriptionRepository();
