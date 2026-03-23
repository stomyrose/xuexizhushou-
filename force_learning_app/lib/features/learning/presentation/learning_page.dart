import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'learning_provider.dart';

class LearningPage extends ConsumerStatefulWidget {
  const LearningPage({super.key});

  @override
  ConsumerState<LearningPage> createState() => _LearningPageState();
}

class _LearningPageState extends ConsumerState<LearningPage> {
  @override
  void initState() {
    super.initState();
    Future.microtask(() {
      ref.read(learningProvider.notifier).loadFiles();
      ref.read(learningProvider.notifier).loadStatistics();
    });
  }

  @override
  Widget build(BuildContext context) {
    final state = ref.watch(learningProvider);

    return Scaffold(
      appBar: AppBar(
        title: const Text('学习中心'),
        actions: [
          IconButton(
            icon: const Icon(Icons.refresh),
            onPressed: () {
              ref.read(learningProvider.notifier).loadFiles();
              ref.read(learningProvider.notifier).loadStatistics();
            },
          ),
        ],
      ),
      body: state.isLoading
          ? const Center(child: CircularProgressIndicator())
          : RefreshIndicator(
              onRefresh: () async {
                await ref.read(learningProvider.notifier).loadFiles();
              },
              child: ListView(
                padding: const EdgeInsets.all(16),
                children: [
                  _buildStatisticsCard(state),
                  const SizedBox(height: 16),
                  _buildRandomLearningCard(),
                  const SizedBox(height: 16),
                  const Text(
                    '知识库',
                    style: TextStyle(
                      fontSize: 18,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  const SizedBox(height: 8),
                  if (state.files.isEmpty)
                    const Card(
                      child: Padding(
                        padding: EdgeInsets.all(32),
                        child: Text(
                          '暂无学习内容',
                          textAlign: TextAlign.center,
                        ),
                      ),
                    )
                  else
                    ...state.files.map((file) => Card(
                      margin: const EdgeInsets.only(bottom: 8),
                      child: ListTile(
                        leading: Container(
                          padding: const EdgeInsets.all(8),
                          decoration: BoxDecoration(
                            color: Colors.blue.withOpacity(0.1),
                            borderRadius: BorderRadius.circular(8),
                          ),
                          child: const Icon(Icons.description, color: Colors.blue),
                        ),
                        title: Text(file.filename),
                        subtitle: Text(file.category ?? '未分类'),
                        trailing: const Icon(Icons.chevron_right),
                        onTap: () {
                          ref.read(learningProvider.notifier).loadRandomFile();
                        },
                      ),
                    )),
                ],
              ),
            ),
    );
  }

  Widget _buildStatisticsCard(LearningState state) {
    final totalSeconds = state.statistics['total_duration_seconds'] ?? 0;
    final hours = totalSeconds ~/ 3600;
    final minutes = (totalSeconds % 3600) ~/ 60;
    final seconds = totalSeconds % 60;

    return Card(
      child: Padding(
        padding: const EdgeInsets.all(20),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const Text(
              '学习统计',
              style: TextStyle(
                fontSize: 16,
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 16),
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceAround,
              children: [
                _buildStatItem('总时长', '${hours}h ${minutes}m'),
                _buildStatItem('学习记录', '${state.records.length}'),
                _buildStatItem('知识文件', '${state.files.length}'),
              ],
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildStatItem(String label, String value) {
    return Column(
      children: [
        Text(
          value,
          style: const TextStyle(
            fontSize: 24,
            fontWeight: FontWeight.bold,
            color: Colors.blue,
          ),
        ),
        const SizedBox(height: 4),
        Text(
          label,
          style: const TextStyle(
            color: Colors.grey,
          ),
        ),
      ],
    );
  }

  Widget _buildRandomLearningCard() {
    return Card(
      color: Colors.blue.shade50,
      child: InkWell(
        onTap: () async {
          await ref.read(learningProvider.notifier).loadRandomFile();
          final file = ref.read(learningProvider).currentFile;
          if (file != null && mounted) {
            _showLearningDialog(file.filename, file.id);
          }
        },
        borderRadius: BorderRadius.circular(12),
        child: const Padding(
          padding: EdgeInsets.all(24),
          child: Row(
            children: [
              Icon(Icons.shuffle, size: 40, color: Colors.blue),
              SizedBox(width: 16),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      '随机学习',
                      style: TextStyle(
                        fontSize: 18,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    SizedBox(height: 4),
                    Text(
                      '点击开始随机学习一个知识点',
                      style: TextStyle(color: Colors.grey),
                    ),
                  ],
                ),
              ),
              Icon(Icons.arrow_forward, color: Colors.blue),
            ],
          ),
        ),
      ),
    );
  }

  void _showLearningDialog(String filename, String fileId) {
    int duration = 0;
    showDialog(
      context: context,
      barrierDismissible: false,
      builder: (context) => StatefulBuilder(
        builder: (context, setState) {
          return AlertDialog(
            title: Text('正在学习: $filename'),
            content: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                Text(
                  '已学习: ${duration}s',
                  style: const TextStyle(fontSize: 24),
                ),
                const SizedBox(height: 16),
                const LinearProgressIndicator(),
              ],
            ),
            actions: [
              TextButton(
                onPressed: () async {
                  final success = await ref.read(learningProvider.notifier).createRecord(
                    fileId: fileId,
                    durationSeconds: duration,
                  );
                  if (context.mounted) {
                    Navigator.pop(context);
                    if (success) {
                      ScaffoldMessenger.of(context).showSnackBar(
                        const SnackBar(content: Text('学习记录已保存')),
                      );
                    }
                  }
                },
                child: const Text('完成学习'),
              ),
            ],
          );
        },
      ),
    );
    Future.doWhile(() async {
      await Future.delayed(const Duration(seconds: 1));
      if (!mounted) return false;
      setState(() => duration++);
      return true;
    });
  }
}
