Mocks contain all auto generated mocks.
to update:

1. delete all folder inside mocks folder, excluding this README.md
2. in terminal run `mockery; Move-Item -Path mocks\internal\* -Destination mocks; Remove-Item mocks\internal`

it will generate all mock for exported interfaces.

To get Mockery, run `go install github.com/vektra/mockery/v2@v2.26.1`
