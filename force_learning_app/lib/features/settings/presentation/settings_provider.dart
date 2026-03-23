import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:shared_preferences/shared_preferences.dart';

class SettingsState {
  final bool popupEnabled;
  final int startHour;
  final int endHour;
  final int durationMinutes;
  final int intervalMinutes;

  SettingsState({
    this.popupEnabled = true,
    this.startHour = 9,
    this.endHour = 22,
    this.durationMinutes = 30,
    this.intervalMinutes = 60,
  });

  SettingsState copyWith({
    bool? popupEnabled,
    int? startHour,
    int? endHour,
    int? durationMinutes,
    int? intervalMinutes,
  }) {
    return SettingsState(
      popupEnabled: popupEnabled ?? this.popupEnabled,
      startHour: startHour ?? this.startHour,
      endHour: endHour ?? this.endHour,
      durationMinutes: durationMinutes ?? this.durationMinutes,
      intervalMinutes: intervalMinutes ?? this.intervalMinutes,
    );
  }
}

class SettingsNotifier extends StateNotifier<SettingsState> {
  SettingsNotifier() : super(SettingsState()) {
    _loadSettings();
  }

  static const String _popupEnabledKey = 'popup_enabled';
  static const String _popupStartHourKey = 'popup_start_hour';
  static const String _popupEndHourKey = 'popup_end_hour';
  static const String _popupDurationKey = 'popup_duration';
  static const String _popupIntervalKey = 'popup_interval';

  Future<void> _loadSettings() async {
    final prefs = await SharedPreferences.getInstance();
    state = SettingsState(
      popupEnabled: prefs.getBool(_popupEnabledKey) ?? true,
      startHour: prefs.getInt(_popupStartHourKey) ?? 9,
      endHour: prefs.getInt(_popupEndHourKey) ?? 22,
      durationMinutes: prefs.getInt(_popupDurationKey) ?? 30,
      intervalMinutes: prefs.getInt(_popupIntervalKey) ?? 60,
    );
  }

  Future<void> setPopupEnabled(bool enabled) async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setBool(_popupEnabledKey, enabled);
    state = state.copyWith(popupEnabled: enabled);
  }

  Future<void> setTimeRange(int startHour, int endHour) async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setInt(_popupStartHourKey, startHour);
    await prefs.setInt(_popupEndHourKey, endHour);
    state = state.copyWith(startHour: startHour, endHour: endHour);
  }

  Future<void> setDuration(int minutes) async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setInt(_popupDurationKey, minutes);
    state = state.copyWith(durationMinutes: minutes);
  }

  Future<void> setInterval(int minutes) async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setInt(_popupIntervalKey, minutes);
    state = state.copyWith(intervalMinutes: minutes);
  }
}

final settingsProvider = StateNotifierProvider<SettingsNotifier, SettingsState>((ref) {
  return SettingsNotifier();
});
