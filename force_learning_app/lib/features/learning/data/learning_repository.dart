import 'package:flutter/foundation.dart';
import '../../../core/api/api_client.dart';
import '../../../core/api/models.dart';

class LearningRepository {
  Future<List<KnowledgeFile>> getFiles({String? category}) async {
    try {
      final response = await apiClient.get(
        '/knowledge/files',
        queryParameters: category != null ? {'category': category} : null,
      );
      
      final apiResponse = ApiResponse.fromJson(response.data);
      if (!apiResponse.isSuccess) {
        throw Exception(apiResponse.message);
      }

      final List<dynamic> data = apiResponse.data ?? [];
      return data.map((e) => KnowledgeFile.fromJson(e)).toList();
    } catch (e) {
      debugPrint('Get files error: $e');
      rethrow;
    }
  }

  Future<KnowledgeFile> getRandomFile() async {
    try {
      final response = await apiClient.get('/knowledge/random');
      final apiResponse = ApiResponse.fromJson(response.data);
      
      if (!apiResponse.isSuccess) {
        throw Exception(apiResponse.message);
      }

      return KnowledgeFile.fromJson(apiResponse.data);
    } catch (e) {
      debugPrint('Get random file error: $e');
      rethrow;
    }
  }

  Future<List<LearningRecord>> getRecords({String? date}) async {
    try {
      final response = await apiClient.get(
        '/learning/records',
        queryParameters: date != null ? {'date': date} : null,
      );
      
      final apiResponse = ApiResponse.fromJson(response.data);
      if (!apiResponse.isSuccess) {
        throw Exception(apiResponse.message);
      }

      final List<dynamic> data = apiResponse.data ?? [];
      return data.map((e) => LearningRecord.fromJson(e)).toList();
    } catch (e) {
      debugPrint('Get records error: $e');
      rethrow;
    }
  }

  Future<LearningRecord> createRecord({
    required String fileId,
    required int durationSeconds,
  }) async {
    try {
      final response = await apiClient.post('/learning/records', data: {
        'file_id': fileId,
        'duration_seconds': durationSeconds,
      });
      
      final apiResponse = ApiResponse.fromJson(response.data);
      if (!apiResponse.isSuccess) {
        throw Exception(apiResponse.message);
      }

      return LearningRecord.fromJson(apiResponse.data);
    } catch (e) {
      debugPrint('Create record error: $e');
      rethrow;
    }
  }

  Future<Map<String, dynamic>> getStatistics() async {
    try {
      final response = await apiClient.get('/learning/statistics');
      final apiResponse = ApiResponse.fromJson(response.data);
      
      if (!apiResponse.isSuccess) {
        throw Exception(apiResponse.message);
      }

      return apiResponse.data ?? {};
    } catch (e) {
      debugPrint('Get statistics error: $e');
      rethrow;
    }
  }
}

final learningRepository = LearningRepository();
