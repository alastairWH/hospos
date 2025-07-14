import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:shared_preferences/shared_preferences.dart';

class ApiService {
  // Get discounts
  static Future<List<Map<String, dynamic>>> getDiscounts() async {
    if (_baseUrl == null) return [];
    try {
      final response = await http.get(Uri.parse('$_baseUrl/discounts'));
      if (response.statusCode == 200) {
        final data = jsonDecode(response.body);
        return List<Map<String, dynamic>>.from(data);
      }
    } catch (_) {}
    return [];
  }
  static String? _baseUrl; // Not set by default

  // Get categories
  static Future<List<String>> getCategories() async {
    if (_baseUrl == null) return [];
    try {
      final response = await http.get(Uri.parse('$_baseUrl/categories'));
      if (response.statusCode == 200) {
        final data = jsonDecode(response.body);
        // Map each object to its 'name' field
        return List<String>.from(data.map((cat) => cat['name']));
      }
    } catch (_) {}
    return [];
  }

  static Future<void> loadServerIp() async {
    final prefs = await SharedPreferences.getInstance();
    final ip = prefs.getString('server_ip');
    if (ip != null && ip.isNotEmpty) {
      _baseUrl = 'http://$ip:8080/api';
    }
  }

  static Future<void> setServerIp(String ip) async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString('server_ip', ip);
    _baseUrl = 'http://$ip:8080/api';
  }

  static String? get baseUrl => _baseUrl;

  static String get currentServerIp => _baseUrl != null ? _baseUrl!.replaceAll(RegExp(r'^https?://|/api/?'), '').replaceAll(':8080', '').replaceAll('/', '') : 'Not set';

  // Terminal linking (graceful error handling, new API)
  static Future<bool> linkTerminal(String code) async {
    if (_baseUrl == null) return false;
    try {
      final response = await http.post(
        Uri.parse('$_baseUrl/linking/link'),
        body: jsonEncode({'linkCode': code, 'deviceInfo': {}}),
        headers: {'Content-Type': 'application/json'},
      ).timeout(const Duration(seconds: 5));
      return response.statusCode == 200;
    } catch (e) {
      // Optionally log error
      return false;
    }
  }

static Future<Map<String, dynamic>?> postSale(Map<String, dynamic> sale) async {
  if (_baseUrl == null) return null;
  try {
    final response = await http.post(
      Uri.parse('$_baseUrl/sales'),
      headers: {'Content-Type': 'application/json'},
      body: jsonEncode(sale),
    );
    if (response.statusCode == 201) {
      return jsonDecode(response.body);
    }
  } catch (_) {}
  return null;
}

  // Heartbeat
  static Future<bool> sendHeartbeat(String tillId, {Map<String, dynamic>? deviceInfo}) async {
    if (_baseUrl == null) return false;
    try {
      final response = await http.post(
        Uri.parse('$_baseUrl/heartbeat'),
        body: jsonEncode({'tillId': tillId, 'deviceInfo': deviceInfo ?? {}}),
        headers: {'Content-Type': 'application/json'},
      );
      return response.statusCode == 200;
    } catch (_) {
      return false;
    }
  }

  // Login
  static Future<Map<String, dynamic>?> login(String userId, String pin) async {
    if (_baseUrl == null) return null;
    try {
      final response = await http.post(
        Uri.parse('$_baseUrl/auth'),
        body: jsonEncode({'name': userId, 'pin': pin}),
        headers: {'Content-Type': 'application/json'},
      );
      if (response.statusCode == 200) {
        return jsonDecode(response.body);
      }
    } catch (_) {}
    return null;
  }

  // Get users
  static Future<List<Map<String, dynamic>>> getUsers() async {
    if (_baseUrl == null) return [];
    try {
      final response = await http.get(Uri.parse('$_baseUrl/users'));
      if (response.statusCode == 200) {
        final data = jsonDecode(response.body);
        // Map _id to id if needed
        return List<Map<String, dynamic>>.from(data).map((user) {
          if (user['id'] == null && user['_id'] != null) {
            user['id'] = user['_id'];
          }
          return user;
        }).toList();
      }
    } catch (_) {}
    return [];
  }

  // Get products
  static Future<List<Map<String, dynamic>>> getProducts() async {
    if (_baseUrl == null) return [];
    try {
      final response = await http.get(Uri.parse('$_baseUrl/products'));
      if (response.statusCode == 200) {
        final data = jsonDecode(response.body);
        return List<Map<String, dynamic>>.from(data);
      }
    } catch (_) {}
    return [];
  }
  // Test connection to server
  static Future<bool> testConnection(String ip, {int port = 8080}) async {
    final url = 'http://$ip:$port/api/heartbeat';
    try {
      final response = await http.post(
        Uri.parse(url),
        body: jsonEncode({'tillId': 'test'}),
        headers: {'Content-Type': 'application/json'},
      ).timeout(const Duration(seconds: 2));
      if (response.statusCode == 200) {
        await setServerIp(ip);
        return true;
      }
    } catch (_) {}
    return false;
  }
}
