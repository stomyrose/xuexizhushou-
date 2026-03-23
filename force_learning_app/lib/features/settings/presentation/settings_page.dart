import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../core/storage/secure_storage.dart';
import 'settings_provider.dart';

class SettingsPage extends ConsumerWidget {
  const SettingsPage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final settings = ref.watch(settingsProvider);

    return Scaffold(
      appBar: AppBar(
        title: const Text('设置'),
      ),
      body: ListView(
        children: [
          const SizedBox(height: 16),
          _buildSection(
            '学习设置',
            [
              SwitchListTile(
                secondary: const Icon(Icons.notifications),
                title: const Text('启用弹窗提醒'),
                subtitle: const Text('定时弹出学习提醒'),
                value: settings.popupEnabled,
                onChanged: (value) {
                  ref.read(settingsProvider.notifier).setPopupEnabled(value);
                },
              ),
              ListTile(
                leading: const Icon(Icons.schedule),
                title: const Text('学习时间段'),
                subtitle: Text('${settings.startHour}:00 - ${settings.endHour}:00'),
                trailing: const Icon(Icons.chevron_right),
                onTap: () => _showTimeRangePicker(context, ref, settings.startHour, settings.endHour),
              ),
              ListTile(
                leading: const Icon(Icons.timer),
                title: const Text('单次学习时长'),
                subtitle: Text('${settings.durationMinutes} 分钟'),
                trailing: const Icon(Icons.chevron_right),
                onTap: () => _showDurationPicker(context, ref, settings.durationMinutes),
              ),
              ListTile(
                leading: const Icon(Icons.repeat),
                title: const Text('弹窗间隔'),
                subtitle: Text('${settings.intervalMinutes} 分钟'),
                trailing: const Icon(Icons.chevron_right),
                onTap: () => _showIntervalPicker(context, ref, settings.intervalMinutes),
              ),
            ],
          ),
          _buildSection(
            '账户设置',
            [
              ListTile(
                leading: const Icon(Icons.sync),
                title: const Text('同步学习记录'),
                subtitle: const Text('将本地记录同步到服务器'),
                trailing: const Icon(Icons.chevron_right),
                onTap: () => _showSyncDialog(context),
              ),
              ListTile(
                leading: const Icon(Icons.lock),
                title: const Text('修改密码'),
                trailing: const Icon(Icons.chevron_right),
                onTap: () {},
              ),
            ],
          ),
          _buildSection(
            '关于',
            [
              ListTile(
                leading: const Icon(Icons.info),
                title: const Text('关于我们'),
                trailing: const Icon(Icons.chevron_right),
                onTap: () => _showAboutDialog(context),
              ),
              ListTile(
                leading: const Icon(Icons.description),
                title: const Text('用户协议'),
                trailing: const Icon(Icons.chevron_right),
                onTap: () {},
              ),
              ListTile(
                leading: const Icon(Icons.privacy_tip),
                title: const Text('隐私政策'),
                trailing: const Icon(Icons.chevron_right),
                onTap: () {},
              ),
            ],
          ),
          const SizedBox(height: 24),
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 16),
            child: OutlinedButton(
              onPressed: () async {
                final confirmed = await showDialog<bool>(
                  context: context,
                  builder: (context) => AlertDialog(
                    title: const Text('确认退出'),
                    content: const Text('确定要退出登录吗？'),
                    actions: [
                      TextButton(
                        onPressed: () => Navigator.pop(context, false),
                        child: const Text('取消'),
                      ),
                      TextButton(
                        onPressed: () => Navigator.pop(context, true),
                        child: const Text('确定'),
                      ),
                    ],
                  ),
                );
                if (confirmed == true && context.mounted) {
                  await SecureStorage.clearAll();
                  Navigator.of(context).pushNamedAndRemoveUntil('/', (route) => false);
                }
              },
              style: OutlinedButton.styleFrom(
                foregroundColor: Colors.red,
                side: const BorderSide(color: Colors.red),
                padding: const EdgeInsets.symmetric(vertical: 16),
              ),
              child: const Text('退出登录'),
            ),
          ),
          const SizedBox(height: 32),
          const Center(
            child: Text(
              'Force Learning v1.0.0',
              style: TextStyle(color: Colors.grey),
            ),
          ),
          const SizedBox(height: 32),
        ],
      ),
    );
  }

  Widget _buildSection(String title, List<Widget> children) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Padding(
          padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
          child: Text(
            title,
            style: const TextStyle(
              fontSize: 14,
              fontWeight: FontWeight.bold,
              color: Colors.grey,
            ),
          ),
        ),
        Card(
          margin: const EdgeInsets.symmetric(horizontal: 16),
          child: Column(children: children),
        ),
        const SizedBox(height: 16),
      ],
    );
  }

  void _showTimeRangePicker(BuildContext context, WidgetRef ref, int startHour, int endHour) {
    int tempStart = startHour;
    int tempEnd = endHour;

    showDialog(
      context: context,
      builder: (context) => StatefulBuilder(
        builder: (context, setState) => AlertDialog(
          title: const Text('设置学习时间段'),
          content: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Row(
                children: [
                  const Text('开始时间: '),
                  const Spacer(),
                  DropdownButton<int>(
                    value: tempStart,
                    items: List.generate(24, (i) => DropdownMenuItem(
                      value: i,
                      child: Text('$i:00'),
                    )),
                    onChanged: (value) {
                      if (value != null) setState(() => tempStart = value);
                    },
                  ),
                ],
              ),
              const SizedBox(height: 16),
              Row(
                children: [
                  const Text('结束时间: '),
                  const Spacer(),
                  DropdownButton<int>(
                    value: tempEnd,
                    items: List.generate(24, (i) => DropdownMenuItem(
                      value: i,
                      child: Text('$i:00'),
                    )),
                    onChanged: (value) {
                      if (value != null) setState(() => tempEnd = value);
                    },
                  ),
                ],
              ),
            ],
          ),
          actions: [
            TextButton(
              onPressed: () => Navigator.pop(context),
              child: const Text('取消'),
            ),
            ElevatedButton(
              onPressed: () {
                if (tempStart < tempEnd) {
                  ref.read(settingsProvider.notifier).setTimeRange(tempStart, tempEnd);
                  Navigator.pop(context);
                } else {
                  ScaffoldMessenger.of(context).showSnackBar(
                    const SnackBar(content: Text('开始时间必须早于结束时间')),
                  );
                }
              },
              child: const Text('确定'),
            ),
          ],
        ),
      ),
    );
  }

  void _showDurationPicker(BuildContext context, WidgetRef ref, int currentDuration) {
    final durations = [15, 20, 30, 45, 60, 90, 120];
    int selected = durations.contains(currentDuration) ? currentDuration : 30;

    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('设置单次学习时长'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: durations.map((d) => RadioListTile<int>(
            title: Text('$d 分钟'),
            value: d,
            groupValue: selected,
            onChanged: (value) {
              if (value != null) {
                selected = value;
                Navigator.pop(context);
                ref.read(settingsProvider.notifier).setDuration(value);
              }
            },
          )).toList(),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('取消'),
          ),
        ],
      ),
    );
  }

  void _showIntervalPicker(BuildContext context, WidgetRef ref, int currentInterval) {
    final intervals = [30, 45, 60, 90, 120, 180];
    int selected = intervals.contains(currentInterval) ? currentInterval : 60;

    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('设置弹窗间隔'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: intervals.map((i) => RadioListTile<int>(
            title: Text('$i 分钟'),
            value: i,
            groupValue: selected,
            onChanged: (value) {
              if (value != null) {
                selected = value;
                Navigator.pop(context);
                ref.read(settingsProvider.notifier).setInterval(value);
              }
            },
          )).toList(),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('取消'),
          ),
        ],
      ),
    );
  }

  void _showSyncDialog(BuildContext context) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('同步学习记录'),
        content: const Text('正在同步学习记录到服务器...'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('关闭'),
          ),
        ],
      ),
    );
  }

  void _showAboutDialog(BuildContext context) {
    showDialog(
      context: context,
      builder: (context) => AboutDialog(
        applicationName: 'Force Learning',
        applicationVersion: '1.0.0',
        applicationIcon: const Icon(Icons.school, size: 48, color: Colors.blue),
        children: const [
          Text('强制学习系统是一款帮助用户通过定时弹窗方式进行强制性学习的工具。'),
        ],
      ),
    );
  }
}
