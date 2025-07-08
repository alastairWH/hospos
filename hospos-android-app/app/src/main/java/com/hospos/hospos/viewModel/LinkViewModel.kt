package com.hospos.hospos.viewModel

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.hospos.hospos.model.LinkResponse
import com.hospos.hospos.repository.HosposRepository
import com.hospos.model.LinkResponse
import com.hospos.repository.HosposRepository
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch

sealed class UiState {
    object Idle : UiState()
    object Loading : UiState()
    data class Success(val response: LinkResponse) : UiState()
    data class Error(val message: String) : UiState()
}

class LinkViewModel : ViewModel() {
    private val repository = HosposRepository()

    private val _uiState = MutableStateFlow<UiState>(UiState.Idle)
    val uiState = _uiState.asStateFlow()

    fun linkTill(code: String, deviceName: String) {
        viewModelScope.launch {
            _uiState.value = UiState.Loading
            try {
                val response = repository.linkTill(code, deviceName)
                _uiState.value = UiState.Success(response)
            } catch (e: Exception) {
                _uiState.value = UiState.Error(e.message ?: "Something went wrong")
            }
        }
    }
}