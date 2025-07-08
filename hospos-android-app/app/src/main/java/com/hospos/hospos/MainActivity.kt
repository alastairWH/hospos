import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Button
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.Icon
import androidx.compose.material3.Text
import androidx.compose.material3.TextField
import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import com.hospos.hospos.viewModel.LinkViewModel
import com.hospos.hospos.viewModel.UiState
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.CheckCircle
import androidx.compose.ui.unit.dp

@Composable
fun LinkingScreen(viewModel: LinkViewModel = viewModel()) {
    val uiState by viewModel.uiState.collectAsState()

    var code by remember { mutableStateOf("") }
    var deviceName by remember { mutableStateOf("Till 1") }
    var status by remember { mutableStateOf<String?>(null) }
    var done by remember { mutableStateOf(false) }

    Column(Modifier.padding(16.dp)) {
        TextField(
            value = code,
            onValueChange = { code = it },
            label = { Text("Enter 12-digit Code") }
        )
        Spacer(Modifier.height(8.dp))
        TextField(
            value = deviceName,
            onValueChange = { deviceName = it },
            label = { Text("Device Name") }
        )
        Spacer(Modifier.height(16.dp))
        Button(onClick = {
            status = "Contacting server…"
            done = false
            viewModel.linkTill(code, deviceName)
        }) {
            Text("Link Device")
        }

        Spacer(Modifier.height(24.dp))
        when (uiState) {
            is UiState.Loading -> {
                status = "Contacting server…"
                Box(Modifier.fillMaxSize(), contentAlignment = Alignment.Center) {
                    CircularProgressIndicator()
                }
                Text(status ?: "")
            }
            is UiState.Success -> {
                status = "Done!"
                done = true
                Icon(Icons.Default.CheckCircle, contentDescription = "Done", tint = Color(0xFF4CAF50))
                Text("Linked! ID: ${(uiState as UiState.Success).response.tillId}")
            }
            is UiState.Error -> {
                status = null
                Text("Error: ${(uiState as UiState.Error).message}", color = Color.Red)
            }
            else -> {}
        }
    }
}

