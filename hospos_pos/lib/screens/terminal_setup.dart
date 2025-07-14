import 'package:flutter/material.dart';
import '../services/api_service.dart';

class TerminalSetupScreen extends StatefulWidget {
  final VoidCallback? onLinked;
  const TerminalSetupScreen({Key? key, this.onLinked}) : super(key: key);

  @override
  State<TerminalSetupScreen> createState() => _TerminalSetupScreenState();
}

class _TerminalSetupScreenState extends State<TerminalSetupScreen> {
  final _terminalIdController = TextEditingController();
  bool _isLinked = false;
  bool _isSetupInProgress = false;
  bool _isOnline = false;
  final _serverIpController = TextEditingController();
  bool _serverConnected = false;
  bool _serverTested = false;
  String? _error;
  String _status = '';

  @override
  void initState() {
    super.initState();
  }


  void _showLinkingCodeDialog() {
    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
        title: const Text('How to get your linking code'),
        content: const Text(
          'To get your linking code, go to the admin dashboard > Terminals > Generate Linking Code. Enter the code here to link this device.'
        ),
        actions: [
          TextButton(
            child: const Text('OK'),
            onPressed: () => Navigator.of(ctx).pop(),
          ),
        ],
      ),
    );
  }

  Future<void> _linkTerminal() async {
    setState(() { _error = null; _isSetupInProgress = true; _status = 'Linking terminal...'; });
    if (_terminalIdController.text.isEmpty) {
      setState(() { _error = 'Please enter a linking code.'; _isSetupInProgress = false; });
      return;
    }
    final success = await ApiService.linkTerminal(_terminalIdController.text);
    if (!success) {
      setState(() { _error = 'Failed to link terminal. Check your code and network.'; _isSetupInProgress = false; });
      return;
    }
    setState(() { _status = 'Checking terminal status...'; });
    final online = await ApiService.sendHeartbeat(_terminalIdController.text);
    setState(() { _isOnline = online; });
    setState(() { _status = online ? 'Terminal is online.' : 'Terminal is offline.'; });
    await Future.delayed(const Duration(seconds: 2));
    setState(() { _status = 'Pulling products...'; });
    await Future.delayed(const Duration(seconds: 2));
    setState(() { _status = 'Finalizing setup...'; });
    await Future.delayed(const Duration(seconds: 2));
    setState(() { _isLinked = true; _isSetupInProgress = false; _status = 'Terminal linked and ready!'; });
    if (widget.onLinked != null) widget.onLinked!();
  }

  @override
  Widget build(BuildContext context) {
    final accentColor = Colors.orangeAccent;
    return Scaffold(
      backgroundColor: Colors.indigo[50],
      appBar: AppBar(
        backgroundColor: Colors.indigo[100],
        elevation: 0,
        title: const Text('Terminal Setup', style: TextStyle(color: Colors.indigo)),
        actions: [
          Row(
            children: [
              if (_serverConnected) ...[
                Icon(Icons.circle, color: Colors.green, size: 18),
                const SizedBox(width: 4),
              ],
              IconButton(
                icon: const Icon(Icons.wifi_tethering, color: Colors.indigo),
                tooltip: 'Test Connection',
                onPressed: () {
                  showDialog(
                    context: context,
                    builder: (ctx) => AlertDialog(
                      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
                      title: const Text('Test Server Connection'),
                      content: Column(
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          TextField(
                            controller: _serverIpController,
                            decoration: InputDecoration(
                              labelText: 'Server IP',
                              border: OutlineInputBorder(borderRadius: BorderRadius.circular(12)),
                              prefixIcon: Icon(Icons.cloud, color: accentColor),
                            ),
                          ),
                          const SizedBox(height: 12),
                          StatefulBuilder(
                            builder: (context, setDialogState) {
                              return ElevatedButton.icon(
                                icon: Icon(Icons.wifi_tethering, color: accentColor),
                                label: const Text('Test Connection'),
                                style: ElevatedButton.styleFrom(
                                  backgroundColor: Colors.indigo,
                                  foregroundColor: Colors.white,
                                  shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                                ),
                                onPressed: () async {
                                  final ip = _serverIpController.text.trim();
                                  if (ip.isEmpty) return;
                                  setDialogState(() { _status = 'Testing connection...'; _serverTested = false; });
                                  final ok = await ApiService.testConnection(ip);
                                  setState(() {
                                    _serverConnected = ok;
                                    _serverTested = true;
                                    _status = ok ? 'Server connected! Using $ip' : 'Server not found.';
                                  });
                                  setDialogState(() {
                                    _serverConnected = ok;
                                    _serverTested = true;
                                    _status = ok ? 'Server connected! Using $ip' : 'Server not found.';
                                  });
                                  if (ok) {
                                    Navigator.of(ctx).pop();
                                  }
                                },
                              );
                            },
                          ),
                          if (_serverTested) ...[
                            const SizedBox(height: 8),
                            Row(
                              mainAxisAlignment: MainAxisAlignment.center,
                              children: [
                                Icon(
                                  _serverConnected ? Icons.circle : Icons.error,
                                  color: _serverConnected ? Colors.green : Colors.red,
                                  size: 18,
                                ),
                                const SizedBox(width: 8),
                                Text(
                                  _serverConnected ? 'Server connected' : 'Server not found',
                                  style: TextStyle(
                                    color: _serverConnected ? Colors.green : Colors.red,
                                    fontWeight: FontWeight.bold,
                                  ),
                                ),
                              ],
                            ),
                            Text(_status, style: TextStyle(color: Colors.indigo[700], fontWeight: FontWeight.w600)),
                          ]
                        ],
                      ),
                    ),
                  );
                },
              ),
              const SizedBox(width: 12),
            ],
          ),
        ],
      ),
      body: Center(
        child: SingleChildScrollView(
          child: Padding(
            padding: const EdgeInsets.all(32.0),
            child: Container(
              decoration: BoxDecoration(
                color: Colors.white,
                borderRadius: BorderRadius.circular(24),
                boxShadow: [
                  BoxShadow(
                    color: Colors.indigo.withOpacity(0.18),
                    blurRadius: 24,
                    offset: const Offset(0, 12),
                  ),
                ],
              ),
              child: Padding(
                padding: const EdgeInsets.all(32.0),
                child: Column(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    Text('Server IP: ${ApiService.baseUrl}', style: TextStyle(color: Colors.indigo[700], fontSize: 14)),
                    const SizedBox(height: 8),
                    Image.asset('assets/logo-hospos.png', height: 80),
                    const SizedBox(height: 24),
                    Text('Terminal Setup', style: TextStyle(
                      fontSize: 28,
                      fontWeight: FontWeight.bold,
                      color: Colors.indigo[700],
                    )),
                    const SizedBox(height: 16),
                    Card(
                      elevation: 0,
                      color: Colors.indigo[50],
                      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
                      child: Padding(
                        padding: const EdgeInsets.all(16.0),
                        child: Column(
                          children: [
                            Text(
                              'To link this terminal, enter the linking code provided by your admin. If you don\'t have a code, tap below for instructions.',
                              style: TextStyle(fontSize: 16, color: Colors.indigo[700]),
                              textAlign: TextAlign.center,
                            ),
                            const SizedBox(height: 12),
                            OutlinedButton.icon(
                              style: OutlinedButton.styleFrom(
                                shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                                side: BorderSide(color: accentColor),
                              ),
                              icon: Icon(Icons.info_outline, color: accentColor),
                              label: const Text('Get Linking Code'),
                              onPressed: _showLinkingCodeDialog,
                            ),
                          ],
                        ),
                      ),
                    ),
                    const SizedBox(height: 24),
                    TextField(
                      controller: _terminalIdController,
                      decoration: InputDecoration(
                        labelText: 'Enter Linking Code',
                        border: OutlineInputBorder(
                          borderRadius: BorderRadius.circular(12),
                        ),
                        prefixIcon: Icon(Icons.vpn_key, color: accentColor),
                      ),
                    ),
                    const SizedBox(height: 20),
                    SizedBox(
                      width: double.infinity,
                      child: ElevatedButton.icon(
                        style: ElevatedButton.styleFrom(
                          backgroundColor: Colors.indigo,
                          foregroundColor: Colors.white,
                          shape: RoundedRectangleBorder(
                            borderRadius: BorderRadius.circular(12),
                          ),
                          elevation: 4,
                          padding: const EdgeInsets.symmetric(vertical: 16),
                        ),
                        icon: Icon(Icons.link, color: accentColor),
                        label: const Text('Link Terminal', style: TextStyle(fontSize: 18)),
                        onPressed: _isSetupInProgress ? null : _linkTerminal,
                      ),
                    ),
                    if (_error != null) ...[
                      const SizedBox(height: 12),
                      Text(_error!, style: const TextStyle(color: Colors.red)),
                    ],
                    if (_isSetupInProgress) ...[
                      const SizedBox(height: 32),
                      Column(
                        children: [
                          const CircularProgressIndicator(),
                          const SizedBox(height: 16),
                          Text(_status, style: TextStyle(
                            color: Colors.indigo[700], fontWeight: FontWeight.w600)),
                        ],
                      ),
                    ],
                    if (_isLinked && !_isSetupInProgress) ...[
                      const SizedBox(height: 32),
                      Row(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          Icon(Icons.check_circle, color: accentColor, size: 28),
                          const SizedBox(width: 8),
                          Text(_status, style: TextStyle(
                            color: Colors.indigo[700], fontWeight: FontWeight.w600)),
                        ],
                      ),
                    ],
                    if (_isOnline) ...[
                      Row(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          Icon(Icons.check_circle, color: Colors.green, size: 24),
                          const SizedBox(width: 8),
                          const Text('Terminal is online', style: TextStyle(color: Colors.green, fontWeight: FontWeight.bold)),
                        ],
                      ),
                    ] else if (_isLinked && !_isOnline) ...[
                      Row(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          Icon(Icons.error, color: Colors.red, size: 24),
                          const SizedBox(width: 8),
                          const Text('Terminal is offline', style: TextStyle(color: Colors.red, fontWeight: FontWeight.bold)),
                        ],
                      ),
                    ],
                  ],
                ),
              ),
            ),
          ),
        ),
      ),
    );
  }
}
