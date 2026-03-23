import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'subscription_provider.dart';

class SubscriptionPage extends ConsumerStatefulWidget {
  const SubscriptionPage({super.key});

  @override
  ConsumerState<SubscriptionPage> createState() => _SubscriptionPageState();
}

class _SubscriptionPageState extends ConsumerState<SubscriptionPage> {
  @override
  void initState() {
    super.initState();
    Future.microtask(() {
      ref.read(subscriptionProvider.notifier).loadPlans();
    });
  }

  @override
  Widget build(BuildContext context) {
    final state = ref.watch(subscriptionProvider);

    return Scaffold(
      appBar: AppBar(
        title: const Text('订阅服务'),
      ),
      body: state.isLoading
          ? const Center(child: CircularProgressIndicator())
          : RefreshIndicator(
              onRefresh: () async {
                await ref.read(subscriptionProvider.notifier).loadPlans();
              },
              child: ListView(
                padding: const EdgeInsets.all(16),
                children: [
                  if (state.currentSubscription != null)
                    Card(
                      color: Colors.green.shade50,
                      child: Padding(
                        padding: const EdgeInsets.all(16),
                        child: Row(
                          children: [
                            const Icon(Icons.check_circle, color: Colors.green, size: 40),
                            const SizedBox(width: 16),
                            Expanded(
                              child: Column(
                                crossAxisAlignment: CrossAxisAlignment.start,
                                children: [
                                  const Text(
                                    '当前订阅',
                                    style: TextStyle(color: Colors.grey),
                                  ),
                                  Text(
                                    state.currentSubscription!.name,
                                    style: const TextStyle(
                                      fontSize: 18,
                                      fontWeight: FontWeight.bold,
                                    ),
                                  ),
                                ],
                              ),
                            ),
                          ],
                        ),
                      ),
                    ),
                  const SizedBox(height: 24),
                  const Text(
                    '选择套餐',
                    style: TextStyle(
                      fontSize: 20,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  const SizedBox(height: 16),
                  ...state.plans.map((plan) => Card(
                    margin: const EdgeInsets.only(bottom: 12),
                    child: InkWell(
                      onTap: state.currentSubscription == null
                          ? () => _handlePurchase(plan.id, plan.name, plan.price)
                          : null,
                      borderRadius: BorderRadius.circular(12),
                      child: Padding(
                        padding: const EdgeInsets.all(20),
                        child: Row(
                          children: [
                            Expanded(
                              child: Column(
                                crossAxisAlignment: CrossAxisAlignment.start,
                                children: [
                                  Text(
                                    plan.name,
                                    style: const TextStyle(
                                      fontSize: 18,
                                      fontWeight: FontWeight.bold,
                                    ),
                                  ),
                                  const SizedBox(height: 4),
                                  Text(
                                    '${plan.durationDays}天',
                                    style: const TextStyle(color: Colors.grey),
                                  ),
                                ],
                              ),
                            ),
                            Column(
                              crossAxisAlignment: CrossAxisAlignment.end,
                              children: [
                                Text(
                                  '¥${plan.price.toStringAsFixed(2)}',
                                  style: const TextStyle(
                                    fontSize: 24,
                                    fontWeight: FontWeight.bold,
                                    color: Colors.blue,
                                  ),
                                ),
                                if (state.currentSubscription == null)
                                  const Text(
                                    '立即订阅',
                                    style: TextStyle(color: Colors.blue),
                                  ),
                              ],
                            ),
                          ],
                        ),
                      ),
                    )),
                  if (state.error != null) ...[
                    const SizedBox(height: 16),
                    Text(
                      state.error!,
                      style: const TextStyle(color: Colors.red),
                      textAlign: TextAlign.center,
                    ),
                  ],
                ],
              ),
            ),
    );
  }

  void _handlePurchase(String planId, String name, double price) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('确认订阅'),
        content: Text('确认订阅 $name，费用 ¥${price.toStringAsFixed(2)}？'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('取消'),
          ),
          ElevatedButton(
            onPressed: () async {
              Navigator.pop(context);
              final success = await ref.read(subscriptionProvider.notifier).purchase(planId);
              if (success && mounted) {
                ScaffoldMessenger.of(context).showSnackBar(
                  const SnackBar(content: Text('订阅成功！')),
                );
              }
            },
            child: const Text('确认'),
          ),
        ],
      ),
    );
  }
}
