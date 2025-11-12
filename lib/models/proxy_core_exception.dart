import 'package:grpc/grpc.dart';





class ProxyCoreException implements Exception {
  
  final Object error;

  
  final StackTrace? stackTrace;

  
  
  
  
  
  const ProxyCoreException(
    this.error, {
    this.stackTrace,
  });

  
  
  
  factory ProxyCoreException.message(String message) {
    return ProxyCoreException(message);
  }

  
  String get message {
    final errorMessage = _extractErrorMessage();
    return errorMessage;
  }

  
  String _extractErrorMessage() {
    if (error is GrpcError) {
      final grpcError = error as GrpcError;
      return _formatGrpcError(grpcError);
    }

    if (error is String) {
      return error as String;
    }

    return error.toString();
  }

  
  String _formatGrpcError(GrpcError error) {
    final code = error.code.toString().replaceAll('StatusCode.', '');
    return 'Code ($code): ${error.message}';
  }

  
  int? get grpcErrorCode {
    if (error is GrpcError) {
      return (error as GrpcError).code;
    }
    return null;
  }

  
  bool get isGrpcError => error is GrpcError;

  @override
  String toString() {
    final baseMessage = message;
    if (stackTrace != null) {
      return '$baseMessage\n$stackTrace';
    }
    return baseMessage;
  }

  @override
  bool operator ==(Object other) {
    if (identical(this, other)) return true;
    return other is ProxyCoreException &&
        other.error.toString() == error.toString();
  }

  @override
  int get hashCode => Object.hash(error, stackTrace);
}
