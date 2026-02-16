use the following command to run the server:

```
bash$ go run main.go run
```

make sure to have the .env and config.yamlfile in the root directory with the following content:

```env
ENVIRONMENT_ENV=development
ENVIRONMENT_PORT=8080
ENVIRONMENT_JWT_SECRET=your_jwt_secret_key
```

config.yaml: as per config.example.yaml
