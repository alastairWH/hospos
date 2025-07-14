import 'package:flutter/material.dart';
import 'screens/terminal_setup.dart';
import 'screens/pin_login.dart';
import 'screens/sales.dart';

void main() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'HOSPOS POS',
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(seedColor: Colors.deepPurple),
      ),
      home: const WorkflowNavigator(),
    );
  }
}

class WorkflowNavigator extends StatefulWidget {
  const WorkflowNavigator({super.key});

  @override
  State<WorkflowNavigator> createState() => _WorkflowNavigatorState();
}

class _WorkflowNavigatorState extends State<WorkflowNavigator> {
  int _step = 0;
  String _userName = 'User';
  String _userRole = 'Role';

  void _nextStep({String? userName, String? userRole}) {
    setState(() {
      _step++;
      if (userName != null) _userName = userName;
      if (userRole != null) _userRole = userRole;
    });
  }

  void _logout() {
    setState(() {
      _step = 1; // Go back to login
    });
  }

  @override
  Widget build(BuildContext context) {
    if (_step == 0) {
      return TerminalSetupScreen(onLinked: _nextStep);
    } else if (_step == 1) {
      return PinLoginScreen(onLogin: () => _nextStep(userName: 'Alastair', userRole: 'Admin'));
    } else {
      return SalesScreen(userName: _userName, userRole: _userRole, onLogout: _logout);
    }
  }
}
