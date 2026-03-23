import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'core/api/api_client.dart';
import 'core/storage/secure_storage.dart';
import 'features/auth/presentation/login_page.dart';
import 'features/auth/presentation/register_page.dart';
import 'features/auth/presentation/home_page.dart';
import 'features/learning/presentation/learning_page.dart';
import 'features/subscription/presentation/subscription_page.dart';
import 'services/popup_service.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  
  await SecureStorage.init();
  await Hive.initFlutter();
  
  runApp(
    const ProviderScope(
      child: ForceLearningApp(),
    ),
  );
}

class ForceLearningApp extends ConsumerStatefulWidget {
  const ForceLearningApp({super.key});

  @override
  ConsumerState<ForceLearningApp> createState() => _ForceLearningAppState();
}

class _ForceLearningAppState extends ConsumerState<ForceLearningApp> {
  late final PopupService _popupService;

  @override
  void initState() {
    super.initState();
    _popupService = PopupService();
    _popupService.initialize();
  }

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Force Learning',
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(seedSeed:  Colors.blue),
        useMaterial3: true,
      ),
      home: const AuthWrapper(),
    );
  }
}

class AuthWrapper extends ConsumerWidget {
  const AuthWrapper({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final authState = ref.watch(authStateProvider);
    
    return authState.when(
      data: (isLoggedIn) => isLoggedIn ? const HomePage() : const LoginPage(),
      loading: () => const Scaffold(
        body: Center(child: CircularProgressIndicator()),
      ),
      error: (_, __) => const LoginPage(),
    );
  }
}

final authStateProvider = FutureProvider<bool>((ref) async {
  final accessToken = await SecureStorage.getAccessToken();
  return accessToken != null && accessToken.isNotEmpty;
});
