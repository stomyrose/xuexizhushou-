class ApiResponse<T> {
  final int code;
  final String message;
  final T? data;

  ApiResponse({
    required this.code,
    required this.message,
    this.data,
  });

  bool get isSuccess => code == 200 || code == 201;

  factory ApiResponse.fromJson(Map<String, dynamic> json) {
    return ApiResponse(
      code: json['code'] ?? 0,
      message: json['message'] ?? '',
      data: json['data'],
    );
  }
}

class AuthResponse {
  final String accessToken;
  final String refreshToken;
  final int expiresIn;

  AuthResponse({
    required this.accessToken,
    required this.refreshToken,
    required this.expiresIn,
  });

  factory AuthResponse.fromJson(Map<String, dynamic> json) {
    return AuthResponse(
      accessToken: json['access_token'] ?? '',
      refreshToken: json['refresh_token'] ?? '',
      expiresIn: json['expires_in'] ?? 0,
    );
  }
}

class UserStatus {
  final String id;
  final String? email;
  final String? phone;
  final int remainingDays;
  final bool isActive;
  final bool hasSubscription;

  UserStatus({
    required this.id,
    this.email,
    this.phone,
    required this.remainingDays,
    required this.isActive,
    required this.hasSubscription,
  });

  factory UserStatus.fromJson(Map<String, dynamic> json) {
    return UserStatus(
      id: json['id'] ?? '',
      email: json['email'],
      phone: json['phone'],
      remainingDays: json['remaining_days'] ?? 0,
      isActive: json['is_active'] ?? false,
      hasSubscription: json['has_subscription'] ?? false,
    );
  }
}

class SubscriptionPlan {
  final String id;
  final String name;
  final int durationDays;
  final double price;
  final bool isActive;

  SubscriptionPlan({
    required this.id,
    required this.name,
    required this.durationDays,
    required this.price,
    required this.isActive,
  });

  factory SubscriptionPlan.fromJson(Map<String, dynamic> json) {
    return SubscriptionPlan(
      id: json['id'] ?? '',
      name: json['name'] ?? '',
      durationDays: json['duration_days'] ?? 0,
      price: (json['price'] ?? 0).toDouble(),
      isActive: json['is_active'] ?? false,
    );
  }
}

class KnowledgeFile {
  final String id;
  final String filename;
  final String filePath;
  final String fileType;
  final String? category;
  final bool isVisible;
  final DateTime uploadedAt;

  KnowledgeFile({
    required this.id,
    required this.filename,
    required this.filePath,
    required this.fileType,
    this.category,
    required this.isVisible,
    required this.uploadedAt,
  });

  factory KnowledgeFile.fromJson(Map<String, dynamic> json) {
    return KnowledgeFile(
      id: json['id'] ?? '',
      filename: json['filename'] ?? '',
      filePath: json['file_path'] ?? '',
      fileType: json['file_type'] ?? '',
      category: json['category'],
      isVisible: json['is_visible'] ?? true,
      uploadedAt: DateTime.tryParse(json['uploaded_at'] ?? '') ?? DateTime.now(),
    );
  }
}

class LearningRecord {
  final String id;
  final String userId;
  final String fileId;
  final DateTime learnedAt;
  final int durationSeconds;

  LearningRecord({
    required this.id,
    required this.userId,
    required this.fileId,
    required this.learnedAt,
    required this.durationSeconds,
  });

  factory LearningRecord.fromJson(Map<String, dynamic> json) {
    return LearningRecord(
      id: json['id'] ?? '',
      userId: json['user_id'] ?? '',
      fileId: json['file_id'] ?? '',
      learnedAt: DateTime.tryParse(json['learned_at'] ?? '') ?? DateTime.now(),
      durationSeconds: json['duration_seconds'] ?? 0,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'file_id': fileId,
      'duration_seconds': durationSeconds,
    };
  }
}
