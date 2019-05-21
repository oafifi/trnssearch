# Transaction Search

### Getting Started

- Clone git repo `git clone https://github.com/oafifi/trnssearch.git`
- Change directory to the project root
- Build Docker image `docker build -t transaction-search .`
- Run Docker container `docker run -p 8081:8081 -it transaction-search`

To run main app `./main`

Endpoint `/api/payment/transaction`

Query Parameters
- provider
- statusCode
- amountMin
- amountMax
- currency

#### Example
`http://localhost:8081/api/payment/transaction?provider=FlypayA&amountMax=500`

### Notes

The fullowing is not completely implemented
- Full test cases
- Specific error messages and response templates
- Input validation

**Filtering transaction items as they are parsed from stream will decrease memory consumption as query result decreases**

**Searching big results can be handled concurrently**
