package com.hospos.hospos.network

import com.hospos.model.*
import retrofit2.http.Body
import retrofit2.http.POST

interface HosposApi {
    @POST("api/linking/link")
    suspend fun linkTill(@Body request: LinkRequest): LinkResponse
}