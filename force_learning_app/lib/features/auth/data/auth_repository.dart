import 'package:flutter/foundation.dart';
import '../../../core/api/api_client.dart';
import '../../../core/api/models.dart';
import '../../../core/storage/secure_storage.dart';

class AuthRepository {
  Future<AuthResponse> register({
    String? email,
    String? phone,
    required String password,
  }) async {
    try {
      final response = await apiClient.post('/auth/register', data: {
        if (email != null) 'email': email,
        if (phone != null) 'phone': phone,
        'password': password,
      });

      final apiResponse = ApiResponse.fromJson(response.data);
      if (!apiResponse.isSuccess) {
        throw Exception(apiResponse.message);
      }

      final authResponse = AuthResponse.fromJson(apiResponse.data);
      await SecureStorage.saveAccessToken(authResponse.accessToken);
      await SecureStorage.saveRefreshToken(authResponse.refreshToken);
      
      return authResponse;
    } catch (e) {
      debugPrint('Register error: $e');
      rethrow;
    }
  }

  Future<AuthResponse> login({
    String? email,
    String? phone,
    required String password,
  }) async {
    try {
      final response = await apiClient.post('/auth/login', data: {
        if (email != null) 'email': email,
        if (phone != null) 'phone': phone,
        'password': password,
      });

      final apiResponse = ApiResponse.fromJson(response.data);
      if (!apiResponse.isSuccess) {
        throw Exception(apiResponse.message);
      }

      final authResponse = AuthResponse.fromJson(apiResponse.data);
      await SecureStorage.saveAccessToken(authResponse.accessToken);
      await SecureStorage.saveRefreshToken(authResponse.refreshToken);
      
      return authResponse;
    } catch (e) {
      debugPrint('Login error: $e');
      rethrow;
    }
  }

  Future<UserStatus> getUserStatus() async {
    try {
      final response = await apiClient.get('/auth/status');
      final apiResponse = ApiResponse.fromJson(response.data);
      
      if (!apiResponse.isSuccess) {
        throw Exception(apiResponse.message);
      }

      return UserStatus.fromJson(apiResponse.data);
    } catch (e) {
      debugPrint('Get status error: $e');
      rethrow;
    }
  }

  Future<void> logout() async {
    await SecureStorage.clearAll();
  }
}

final authRepository = AuthRepository();
