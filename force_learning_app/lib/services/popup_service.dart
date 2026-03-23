import 'dart:async';
import 'package:flutter/material.dart';
import 'package:shared_preferences/shared_preferences.dart';

class PopupService {
  static const String _popupEnabledKey = 'popup_enabled';
  static const String _popupStartHourKey = 'popup_start_hour';
  static const String _popupEndHourKey = 'popup_end_hour';
  static const String _popupDurationKey = 'popup_duration';
  static const String _popupIntervalKey = 'popup_interval';

  Timer? _popupTimer;
  bool _isPopupShowing = false;
  DateTime? _lastPopupTime;
  VoidCallback? _onPopupShown;
  VoidCallback? _onPopupDismissed;
  Function(String fileId)? _onStartLearning;

  bool _popupEnabled = true;
  int _startHour = 9;
  int _endHour = 22;
  int _durationMinutes = 30;
  int _intervalMinutes = 60;

  bool get isPopupShowing => _isPopupShowing;
  bool get popupEnabled => _popupEnabled;
  int get startHour => _startHour;
  int get endHour => _endHour;
  int get durationMinutes => _durationMinutes;
  int get intervalMinutes => _intervalMinutes;

  Future<void> initialize() async {
    await _loadSettings();
    _startTimer();
  }

  Future<void> _loadSettings() async {
    try {
      final prefs = await SharedPreferences.getInstance();
      _popupEnabled = prefs.getBool(_popupEnabledKey) ?? true;
      _startHour = prefs.getInt(_popupStartHourKey) ?? 9;
      _endHour = prefs.getInt(_popupEndHourKey) ?? 22;
      _durationMinutes = prefs.getInt(_popupDurationKey) ?? 30;
      _intervalMinutes = prefs.getInt(_popupIntervalKey) ?? 60;
    } catch (e) {
      debugPrint('Failed to load popup settings: $e');
    }
  }

  void _startTimer() {
    _popupTimer?.cancel();
    _popupTimer = Timer.periodic(const Duration(minutes: 1), (_) {
      _checkAndShowPopup();
    });
  }

  void _checkAndShowPopup() {
    if (!_popupEnabled) return;

    final now = DateTime.now();
    if (now.hour < _startHour || now.hour >= _endHour) return;

    if (_lastPopupTime != null) {
      final minutesSinceLastPopup = now.difference(_lastPopupTime!).inMinutes;
      if (minutesSinceLastPopup < _intervalMinutes) return;
    }

    _showLearningPopup();
  }

  void _showLearningPopup() {
    if (_isPopupShowing) return;

    _isPopupShowing = true;
    _lastPopupTime = DateTime.now();
    _onPopupShown?.call();

    debugPrint('Learning popup shown at ${DateTime.now()}');
  }

  void dismissPopup() {
    _isPopupShowing = false;
    _onPopupDismissed?.call();
  }

  Future<void> setPopupEnabled(bool enabled) async {
    _popupEnabled = enabled;
    await _saveSettings();
  }

  Future<void> setTimeRange(int startHour, int endHour) async {
    _startHour = startHour;
    _endHour = endHour;
    await _saveSettings();
  }

  Future<void> setDuration(int minutes) async {
    _durationMinutes = minutes;
    await _saveSettings();
  }

  Future<void> setInterval(int minutes) async {
    _intervalMinutes = minutes;
    await _saveSettings();
  }

  Future<void> _saveSettings() async {
    try {
      final prefs = await SharedPreferences.getInstance();
      await prefs.setBool(_popupEnabledKey, _popupEnabled);
      await prefs.setInt(_popupStartHourKey, _startHour);
      await prefs.setInt(_popupEndHourKey, _endHour);
      await prefs.setInt(_popupDurationKey, _durationMinutes);
      await prefs.setInt(_popupIntervalKey, _intervalMinutes);
    } catch (e) {
      debugPrint('Failed to save popup settings: $e');
    }
  }

  void setOnPopupShown(VoidCallback callback) {
    _onPopupShown = callback;
  }

  void setOnPopupDismissed(VoidCallback callback) {
    _onPopupDismissed = callback;
  }

  void setOnStartLearning(Function(String fileId) callback) {
    _onStartLearning = callback;
  }

  void dispose() {
    _popupTimer?.cancel();
  }
}

class LearningPopupDialog extends StatefulWidget {
  final String title;
  final String message;
  final int durationMinutes;
  final VoidCallback onStart;
  final VoidCallback onDismiss;

  const LearningPopupDialog({
    super.key,
    required this.title,
    required this.message,
    required this.durationMinutes,
    required this.onStart,
    required this.onDismiss,
  });

  @override
  State<LearningPopupDialog> createState() => _LearningPopupDialogState();
}

class _LearningPopupDialogState extends State<LearningPopupDialog> {
  late int _remainingSeconds;
  late Timer _timer;
  bool _isLearning = false;

  @override
  void initState() {
    super.initState();
    _remainingSeconds = widget.durationMinutes * 60;
    _startTimer();
  }

  void _startTimer() {
    _timer = Timer.periodic(const Duration(seconds: 1), (timer) {
      if (_remainingSeconds > 0) {
        setState(() {
          _remainingSeconds--;
        });
      } else {
        timer.cancel();
        _completeLearning();
      }
    });
  }

  void _completeLearning() {
    setState(() {
      _isLearning = true;
    });
    widget.onStart();
  }

  @override
  void dispose() {
    _timer.cancel();
    super.dispose();
  }

  String _formatTime(int seconds) {
    final minutes = seconds ~/ 60;
    final secs = seconds % 60;
    return '${minutes.toString().padLeft(2, '0')}:${secs.toString().padLeft(2, '0')}';
  }

  @override
  Widget build(BuildContext context) {
    final progress = 1 - (_remainingSeconds / (widget.durationMinutes * 60));

    return Dialog(
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(20),
      ),
      child: Padding(
        padding: const EdgeInsets.all(24),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Container(
              padding: const EdgeInsets.all(16),
              decoration: BoxDecoration(
                color: Colors.blue.shade50,
                shape: BoxShape.circle,
              ),
              child: Icon(
                _isLearning ? Icons.check_circle : Icons.school,
                size: 48,
                color: Colors.blue,
              ),
            ),
            const SizedBox(height: 20),
            Text(
              _isLearning ? '学习完成' : widget.title,
              style: const TextStyle(
                fontSize: 22,
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 8),
            Text(
              _isLearning 
                  ? '恭喜完成本次学习！' 
                  : widget.message,
              textAlign: TextAlign.center,
              style: TextStyle(
                fontSize: 14,
                color: Colors.grey[600],
              ),
            ),
            const SizedBox(height: 24),
            if (!_isLearning) ...[
              Text(
                _formatTime(_remainingSeconds),
                style: const TextStyle(
                  fontSize: 48,
                  fontWeight: FontWeight.bold,
                  fontFamily: 'monospace',
                  color: Colors.blue,
                ),
              ),
              const SizedBox(height: 8),
              LinearProgressIndicator(
                value: progress,
                backgroundColor: Colors.grey[200],
                valueColor: AlwaysStoppedAnimation<Color>(
                  progress > 0.7 ? Colors.green : Colors.blue,
                ),
              ),
              const SizedBox(height: 8),
              Text(
                '学习进度 ${(progress * 100).toInt()}%',
                style: TextStyle(
                  color: Colors.grey[600],
                  fontSize: 12,
                ),
              ),
            ] else ...[
              const Icon(
                Icons.celebration,
                size: 48,
                color: Colors.amber,
              ),
              const SizedBox(height: 8),
              Text(
                '已学习 ${widget.durationMinutes} 分钟',
                style: const TextStyle(
                  fontSize: 16,
                  color: Colors.green,
                ),
              ),
            ],
            const SizedBox(height: 24),
            Row(
              children: [
                Expanded(
                  child: OutlinedButton(
                    onPressed: () {
                      widget.onDismiss();
                      Navigator.of(context).pop();
                    },
                    style: OutlinedButton.styleFrom(
                      padding: const EdgeInsets.symmetric(vertical: 12),
                    ),
                    child: const Text('跳过'),
                  ),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: ElevatedButton(
                    onPressed: _isLearning
                        ? () {
                            widget.onDismiss();
                            Navigator.of(context).pop();
                          }
                        : null,
                    style: ElevatedButton.styleFrom(
                      backgroundColor: Colors.blue,
                      padding: const EdgeInsets.symmetric(vertical: 12),
                    ),
                    child: Text(_isLearning ? '完成' : '学习中...'),
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}
