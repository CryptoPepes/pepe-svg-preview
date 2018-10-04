# Genetic testing

Small development server to test the Genetic expression code.
Pepe DNA is semi-randomly generated,
 i.e. one side of the DNA fully random,
 the other side replicates the first one with some variation
 (exactly how the contract generates new Gen0 pepes).

## Setup

A unix system is highly preferred, for better troubleshooting help.

1) Install Go (Lookup a tutorial)
1) Add GOPATH to your Path (Should be included in any Go install tutorial)
1) Run the following commands:

```bash
# Install cryptopepe projects
cd "$GOPATH/src/"
mkdir "cryptopepe.io"
cd "cryptopepe.io"
git clone https://github.com/CryptoPepes/pepe-svg-preview.git cryptopepe-svg-preview
git clone https://github.com/CryptoPepes/pepe-svg-build.git cryptopepe-svg
git clone https://github.com/CryptoPepes/pepe-reader.git cryptopepe-reader

# install dependencies (MUX for routing)
cd cryptopepe-svg-preview
go get ./...

govendor add +e
# Explicitly add resource files
govendor add cryptopepe.io/cryptopepe-svg/builder/tmpl/^
```

## Running

```bash
# To start the program: (Note, this works, but using an IDE may be a better idea)
# from the cryptopepe-svg-preview project root:
go run main.go
```

## Changing

Gene expressors can be found here: `$GOPATH/cryptopepe.io/cryptopepe-reader/pepe/expressors.go`.
Simply restart the server (should be *very* fast) to view the results of the new gene expressor code.

To push changes to the repository:

```bash
# From "$GOPATH/cryptopepe.io/cryptopepe-reader"
# Create a new branch
git checkout -b gene-changes
# Add changes to git stage
git add .
# commit changes, add descriptive message
git commit -m "Made changes to gene expressor"
# push to repo!
git push origin gene-changes
```
