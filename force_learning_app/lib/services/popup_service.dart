import 'dart:async';
import 'package:flutter/foundation.dart';
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

  bool _popupEnabled = true;
  int _startHour = 9;
  int _endHour = 22;
  int _durationMinutes = 30;
  int _intervalMinutes = 60;

  bool get isPopupShowing => _isPopupShowing;

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
  }

  void setPopupEnabled(bool enabled) {
    _popupEnabled = enabled;
    _saveSettings();
  }

  void setTimeRange(int startHour, int endHour) {
    _startHour = startHour;
    _endHour = endHour;
    _saveSettings();
  }

  void setDuration(int minutes) {
    _durationMinutes = minutes;
    _saveSettings();
  }

  void setInterval(int minutes) {
    _intervalMinutes = minutes;
    _saveSettings();
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

  void dispose() {
    _popupTimer?.cancel();
  }
}
