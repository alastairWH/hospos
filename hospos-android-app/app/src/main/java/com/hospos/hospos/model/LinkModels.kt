package com.hospos.hospos.model

data class LinkRequest(
    val linkCode: String,
    val deviceInfo: DeviceInfo
)

data class DeviceInfo(
    val platform: String = "android",
    val deviceName: String
)

data class LinkResponse(
    val success: Boolean,
    val tillId: String,
    val initialData: InitialData
)

data class InitialData(
    val products: List<Product>,
    val categories: List<Category>,
    val users: List<User>,
    val roles: List<Role>
)

data class Product(val id: String, val name: String, val price: Double)
data class Category(val id: String, val name: String)
data class User(val id: String, val username: String)
data class Role(val id: String, val name: String)