# Instructional Context: `ets-gruppenadressgenerator`

This file provides the operational context for Gemini. For detailed business logic, requirements, and project scope, always refer to the **Source of Truth** in `spec/01-main.md`.

## Key Files & Structure
- **`spec/01-main.md`**: **Core Specification.** Contains all business requirements, naming conventions, and the KNX hierarchy logic.
- **`main.go`**: The entry point and main logic of the Golang application.
- **`input.json`**: Current project configuration (rooms, trades, devices).
- **`knx-object-types.json`**: Mapping of device types to KNX functions.
- **`spec/group-addresses-example1.xml`**: Reference for the expected XML output format.

## Technical Stack
- **Language:** Go (Golang)
- **Modul:** `github.com/maxarndt/ets-gruppenadressgenerator`
- **Output:** XML (`output.xml`) in KNX `GroupAddress-Export` format.

## Building and Running
- **Run:** `go run main.go`
- **Build:** `go build -o ets-gen main.go`
- **Test:** `go test ./...`

## Development Guidelines
- **Tests:** Always update or add tests when making code changes to ensure continued reliability.
- **Logic:** Adhere strictly to the rules defined in `spec/01-main.md`.
- **Error Handling:** Terminate execution if a device type is not found in `knx-object-types.json`.
- **Conventions:** Follow standard Go idioms. Use `output.xml` as the default output file name.
