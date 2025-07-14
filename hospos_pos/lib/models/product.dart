class Product {
  final String id;
  final String name;
  final double price;
  final String category;

  Product({required this.id, required this.name, required this.price, required this.category});

  factory Product.fromJson(Map<String, dynamic> json) {
    return Product(
      id: json['id'] ?? '',
      name: json['name'] ?? '',
      price: (json['price'] ?? 0).toDouble(),
      category: json['category'] ?? '',
    );
  }
}
