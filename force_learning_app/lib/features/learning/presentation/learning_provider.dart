import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../data/learning_repository.dart';
import '../../../core/api/models.dart';

class LearningState {
  final List<KnowledgeFile> files;
  final List<LearningRecord> records;
  final KnowledgeFile? currentFile;
  final Map<String, dynamic> statistics;
  final bool isLoading;
  final String? error;

  LearningState({
    this.files = const [],
    this.records = const [],
    this.currentFile,
    this.statistics = const {},
    this.isLoading = false,
    this.error,
  });

  LearningState copyWith({
    List<KnowledgeFile>? files,
    List<LearningRecord>? records,
    KnowledgeFile? currentFile,
    Map<String, dynamic>? statistics,
    bool? isLoading,
    String? error,
  }) {
    return LearningState(
      files: files ?? this.files,
      records: records ?? this.records,
      currentFile: currentFile ?? this.currentFile,
      statistics: statistics ?? this.statistics,
      isLoading: isLoading ?? this.isLoading,
      error: error,
    );
  }
}

class LearningNotifier extends StateNotifier<LearningState> {
  LearningNotifier() : super(LearningState());

  Future<void> loadFiles() async {
    state = state.copyWith(isLoading: true, error: null);
    try {
      final files = await learningRepository.getFiles();
      state = state.copyWith(files: files, isLoading: false);
    } catch (e) {
      state = state.copyWith(isLoading: false, error: e.toString());
    }
  }

  Future<void> loadRandomFile() async {
    state = state.copyWith(isLoading: true, error: null);
    try {
      final file = await learningRepository.getRandomFile();
      state = state.copyWith(currentFile: file, isLoading: false);
    } catch (e) {
      state = state.copyWith(isLoading: false, error: e.toString());
    }
  }

  Future<void> loadRecords() async {
    try {
      final records = await learningRepository.getRecords();
      state = state.copyWith(records: records);
    } catch (e) {
      debugPrint('Load records error: $e');
    }
  }

  Future<void> loadStatistics() async {
    try {
      final statistics = await learningRepository.getStatistics();
      state = state.copyWith(statistics: statistics);
    } catch (e) {
      debugPrint('Load statistics error: $e');
    }
  }

  Future<bool> createRecord({
    required String fileId,
    required int durationSeconds,
  }) async {
    try {
      await learningRepository.createRecord(
        fileId: fileId,
        durationSeconds: durationSeconds,
      );
      await loadRecords();
      await loadStatistics();
      return true;
    } catch (e) {
      state = state.copyWith(error: e.toString());
      return false;
    }
  }
}

final learningProvider = StateNotifierProvider<LearningNotifier, LearningState>((ref) {
  return LearningNotifier();
});
