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

  void _nextStep() {
    setState(() { _step++; });
  }

  @override
  Widget build(BuildContext context) {
    if (_step == 0) {
      return TerminalSetupScreen(onLinked: _nextStep);
    } else if (_step == 1) {
      return PinLoginScreen(onLogin: _nextStep);
    } else {
      return const SalesScreen();
    }
  }
}
