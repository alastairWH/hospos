package com.hospos.hospos.repository

import com.hospos.model.LinkRequest
import com.hospos.model.LinkResponse
import com.hospos.network.ApiService

class HosposRepository {
    suspend fun linkTill(code: String, deviceName: String): LinkResponse {
        val request = LinkRequest(
            linkCode = code,
            deviceInfo = com.hospos.model.DeviceInfo(deviceName = deviceName)
        )
        return ApiService.api.linkTill(request)
    }
}