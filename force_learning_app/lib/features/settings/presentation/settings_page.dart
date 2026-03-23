import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../core/storage/secure_storage.dart';

class SettingsPage extends ConsumerWidget {
  const SettingsPage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
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
              ListTile(
                leading: const Icon(Icons.timer),
                title: const Text('学习提醒时间'),
                subtitle: const Text('每天 09:00'),
                trailing: const Icon(Icons.chevron_right),
                onTap: () {},
              ),
              ListTile(
                leading: const Icon(Icons.schedule),
                title: const Text('单次学习时长'),
                subtitle: const Text('30 分钟'),
                trailing: const Icon(Icons.chevron_right),
                onTap: () {},
              ),
              SwitchListTile(
                secondary: const Icon(Icons.notifications),
                title: const Text('启用弹窗提醒'),
                value: true,
                onChanged: (value) {},
              ),
            ],
          ),
          _buildSection(
            '账户设置',
            [
              ListTile(
                leading: const Icon(Icons.person),
                title: const Text('修改个人信息'),
                trailing: const Icon(Icons.chevron_right),
                onTap: () {},
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
                onTap: () {},
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
                await SecureStorage.clearAll();
                if (context.mounted) {
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
}
