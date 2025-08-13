# user-go (OTP TDD Starter)

This repository is a TDD-first scaffold for an OTP-based auth service in Go.

Steps to start (local dev):

1. Initialize modules and download deps:

```bash
# from project root
go mod tidy
```

2. Run unit tests (they are expected to fail at first — RED step of TDD):

```bash
go test ./... -v
```

3. Next: implement `OtpService.RequestOTP` so tests turn GREEN.
