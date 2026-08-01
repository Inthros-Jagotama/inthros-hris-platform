// licensectl — CLI utility untuk mesin lisensi On-Premise.
//
// Usage:
//
//	licensectl gen-key [--bits 2048] [--out private.pem] [--pub public.pem]
//	licensectl gen-lic --priv private.pem [--out license.lic] \
//	    --company-id <uuid> --company "<name>" \
//	    --expires 2027-12-31 --modules organization,employee,payroll \
//	    --max-employees 500
//
// Output .lic berisi JSON {payload, signature} (RSA-SHA256) yang diverifikasi
// backend saat DEPLOYMENT_MODE=ON_PREMISE.
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/inthros/hris-platform/internal/pkg/onpremise"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "gen-key":
		runGenKey(os.Args[2:])
	case "gen-lic":
		runGenLic(os.Args[2:])
	case "help", "-h", "--help":
		usage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", os.Args[1])
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Println(`licensectl — On-Premise license utility

Commands:
  gen-key   Generate RSA keypair (private + public PEM)
  gen-lic   Generate .lic license file signed with the private key

Examples:
  licensectl gen-key --bits 2048 --out private.pem --pub public.pem
  licensectl gen-lic --priv private.pem --out license.lic \
      --company-id 00000000-0000-0000-0000-000000000001 \
      --company "PT Contoh" --expires 2027-12-31 \
      --modules organization,employee,payroll --max-employees 500`)
}

func runGenKey(args []string) {
	fs := flag.NewFlagSet("gen-key", flag.ExitOnError)
	bits := fs.Int("bits", onpremise.DefaultKeySize, "RSA key size (bits)")
	out := fs.String("out", "private.pem", "output path for private key")
	pub := fs.String("pub", "public.pem", "output path for public key")
	if err := fs.Parse(args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	privPEM, pubPEM, err := onpremise.GenerateKeyPair(*bits)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to generate keypair: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(*out, privPEM, 0o600); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write private key: %v\n", err)
		os.Exit(1)
	}
	if err := os.WriteFile(*pub, pubPEM, 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write public key: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Keypair generated:\n  private: %s\n  public:  %s\n  bits:    %d\n", *out, *pub, *bits)
}

func runGenLic(args []string) {
	fs := flag.NewFlagSet("gen-lic", flag.ExitOnError)
	priv := fs.String("priv", "", "path to RSA private key PEM (required)")
	out := fs.String("out", "license.lic", "output path for .lic file")
	companyID := fs.String("company-id", "", "company UUID (required)")
	companyName := fs.String("company", "", "company name (required)")
	expires := fs.String("expires", "", "expiry date YYYY-MM-DD (required)")
	modules := fs.String("modules", "", "comma-separated allowed module slugs (required)")
	maxEmployees := fs.Int("max-employees", 0, "max employees allowed (required)")
	if err := fs.Parse(args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if *priv == "" || *companyID == "" || *companyName == "" || *expires == "" || *modules == "" || *maxEmployees <= 0 {
		fmt.Fprintln(os.Stderr, "Missing required flag(s). Use --help for usage.")
		os.Exit(1)
	}

	expiry, err := time.Parse("2006-01-02", *expires)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid --expires (use YYYY-MM-DD): %v\n", err)
		os.Exit(1)
	}

	privPEM, err := os.ReadFile(*priv)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read private key: %v\n", err)
		os.Exit(1)
	}

	lic := onpremise.License{
		CompanyID:      *companyID,
		CompanyName:    *companyName,
		ExpiresAt:      expiry,
		AllowedModules: splitCSV(*modules),
		MaxEmployees:   *maxEmployees,
	}

	if err := onpremise.WriteLicenseFile(*out, privPEM, lic); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write license: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("License written: %s\n  company: %s (%s)\n  expires: %s\n  modules: %v\n  max_employees: %d\n",
		*out, lic.CompanyName, lic.CompanyID, lic.ExpiresAt.Format("2006-01-02"), lic.AllowedModules, lic.MaxEmployees)
}

func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}
