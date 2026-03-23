import 'package:hive/hive.dart';
import '../core/api/api_client.dart';
import '../core/storage/secure_storage.dart';

class SyncService {
  static const String _recordsBoxName = 'learning_records';
  static const String _lastSyncKey = 'last_sync_time';

  Box? _recordsBox;
  DateTime? _lastSyncTime;

  Future<void> init() async {
    _recordsBox = await Hive.openBox(_recordsBoxName);
    final lastSyncStr = _recordsBox?.get(_lastSyncKey);
    if (lastSyncStr != null) {
      _lastSyncTime = DateTime.tryParse(lastSyncStr);
    }
  }

  Future<void> saveLocalRecord(LearningRecordDTO record) async {
    await _recordsBox?.put(record.clientId, record.toJson());
  }

  Future<List<LearningRecordDTO>> getLocalRecords() async {
    final records = <LearningRecordDTO>[];
    _recordsBox?.toMap().forEach((key, value) {
      if (value is Map) {
        records.add(LearningRecordDTO.fromJson(Map<String, dynamic>.from(value)));
      }
    });
    return records;
  }

  Future<SyncResult> syncRecords() async {
    if (_recordsBox == null) {
      await init();
    }

    final localRecords = await getLocalRecords();
    if (localRecords.isEmpty) {
      return SyncResult(
        success: true,
        syncedCount: 0,
        serverTime: DateTime.now().toIso8601String(),
      );
    }

    try {
      final response = await apiClient.post(
        '/learning/sync',
        data: {
          'last_sync_time': _lastSyncTime?.toIso8601String() ?? '',
          'records': localRecords.map((r) => r.toSyncJson()).toList(),
        },
      );

      if (response.statusCode == 200 && response.data['code'] == 200) {
        final data = response.data['data'];
        final syncedRecords = data['synced_records'] as List? ?? [];
        final serverTime = data['server_time'] as String? ?? DateTime.now().toIso8601String();

        for (final synced in syncedRecords) {
          if (synced['synced'] == true) {
            await _recordsBox?.delete(synced['client_id']);
          }
        }

        _lastSyncTime = DateTime.now();
        await _recordsBox?.put(_lastSyncKey, _lastSyncTime!.toIso8601String());

        return SyncResult(
          success: true,
          syncedCount: syncedRecords.where((s) => s['synced'] == true).length,
          serverTime: serverTime,
        );
      } else {
        return SyncResult(
          success: false,
          error: response.data['message'] ?? 'Sync failed',
          serverTime: DateTime.now().toIso8601String(),
        );
      }
    } catch (e) {
      return SyncResult(
        success: false,
        error: e.toString(),
        serverTime: DateTime.now().toIso8601String(),
      );
    }
  }

  Future<List<LearningRecordDTO>> fetchUnsyncedRecords() async {
    try {
      final response = await apiClient.get(
        '/learning/unsynced',
        queryParameters: {
          'since': _lastSyncTime?.toIso8601String() ?? '',
        },
      );

      if (response.statusCode == 200 && response.data['code'] == 200) {
        final records = response.data['data'] as List? ?? [];
        return records.map((r) => LearningRecordDTO.fromServerJson(r)).toList();
      }
    } catch (e) {
      debugPrint('Failed to fetch unsynced records: $e');
    }
    return [];
  }

  DateTime? get lastSyncTime => _lastSyncTime;
}

class LearningRecordDTO {
  final String clientId;
  final String fileId;
  final int durationSeconds;
  final DateTime learnedAt;
  bool synced;

  LearningRecordDTO({
    required this.clientId,
    required this.fileId,
    required this.durationSeconds,
    required this.learnedAt,
    this.synced = false,
  });

  factory LearningRecordDTO.fromJson(Map<String, dynamic> json) {
    return LearningRecordDTO(
      clientId: json['client_id'] ?? '',
      fileId: json['file_id'] ?? '',
      durationSeconds: json['duration_seconds'] ?? 0,
      learnedAt: DateTime.tryParse(json['learned_at'] ?? '') ?? DateTime.now(),
      synced: json['synced'] ?? false,
    );
  }

  factory LearningRecordDTO.fromServerJson(Map<String, dynamic> json) {
    return LearningRecordDTO(
      clientId: json['id'] ?? '',
      fileId: json['file_id'] ?? '',
      durationSeconds: json['duration_seconds'] ?? 0,
      learnedAt: DateTime.tryParse(json['learned_at'] ?? '') ?? DateTime.now(),
      synced: true,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'client_id': clientId,
      'file_id': fileId,
      'duration_seconds': durationSeconds,
      'learned_at': learnedAt.toIso8601String(),
      'synced': synced,
    };
  }

  Map<String, dynamic> toSyncJson() {
    return {
      'client_id': clientId,
      'file_id': fileId,
      'duration_seconds': durationSeconds,
      'learned_at': learnedAt.toIso8601String(),
    };
  }
}

class SyncResult {
  final bool success;
  final int syncedCount;
  final String serverTime;
  final String? error;

  SyncResult({
    required this.success,
    required this.syncedCount,
    required this.serverTime,
    this.error,
  });
}

final syncService = SyncService();
