# Savvy Shopper

A CLI tool to compare product prices from Amazon and Walmart using ZincAPI. Get the best deals by searching both retailers simultaneously and comparing prices side by side.

## Features

- 🔍 Search products on Amazon and Walmart simultaneously
- 💰 Compare prices across retailers
- 🚀 Fast concurrent searches with retry/backoff
- 📊 Clean tabular output with product details
- 🔑 Secure API key management

## Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/savvyshopper.git
cd savvyshopper

# Build the binary
go build -o savvyshopper ./cmd/main.go

# Move to your PATH (optional)
mv savvyshopper /usr/local/bin/
```

## Configuration

Set your ZincAPI key as an environment variable:

```bash
export ZINC_API_KEY="your-api-key-here"
```

## Usage

### Basic Usage

```bash
# Search for a product
savvyshopper "AirPods Pro 2nd Gen"

# Or run directly with go
go run ./cmd/main.go "AirPods Pro 2nd Gen"
```

### Interactive Mode

If no product name is provided, you'll be prompted to enter one:

```bash
$ savvyshopper
Enter product: AirPods Pro 2nd Gen
```

### Example Output

```
Product Title                                    Price    Retailer    URL
AirPods Pro (2nd Generation)                    $199.99  Amazon      https://amazon.com/...
AirPods Pro (2nd Generation) with MagSafe Case  $249.99  Walmart     https://walmart.com/...
```

## Development

```bash
# Run tests
make test

# Run linter
make lint

# Run the CLI
make run
```

## Error Handling

The tool handles various error cases gracefully:

- 🔒 Missing or invalid API key
- 🌐 Network connectivity issues
- 🔍 No results found
- ⚠️ Invalid product names

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
