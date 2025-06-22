# fc-pos-go-labs-deploy-cloudrun

## Descrição do desafio
```
Objetivo: Desenvolver um sistema em Go que receba um CEP, identifica a cidade e retorna o clima atual (temperatura em graus celsius, fahrenheit e kelvin). Esse sistema deverá ser publicado no Google Cloud Run.

Requisitos:

O sistema deve receber um CEP válido de 8 digitos
O sistema deve realizar a pesquisa do CEP e encontrar o nome da localização, a partir disso, deverá retornar as temperaturas e formata-lás em: Celsius, Fahrenheit, Kelvin.
O sistema deve responder adequadamente nos seguintes cenários:
Em caso de sucesso:
Código HTTP: 200
Response Body: { "temp_C": 28.5, "temp_F": 28.5, "temp_K": 28.5 }
Em caso de falha, caso o CEP não seja válido (com formato correto):
Código HTTP: 422
Mensagem: invalid zipcode
​​​Em caso de falha, caso o CEP não seja encontrado:
Código HTTP: 404
Mensagem: can not find zipcode
Deverá ser realizado o deploy no Google Cloud Run.

Dicas:

Utilize a API viaCEP (ou similar) para encontrar a localização que deseja consultar a temperatura: https://viacep.com.br/
Utilize a API WeatherAPI (ou similar) para consultar as temperaturas desejadas: https://www.weatherapi.com/
Para realizar a conversão de Celsius para Fahrenheit, utilize a seguinte fórmula: F = C * 1,8 + 32
Para realizar a conversão de Celsius para Kelvin, utilize a seguinte fórmula: K = C + 273
Sendo F = Fahrenheit
Sendo C = Celsius
Sendo K = Kelvin

Entrega:

O código-fonte completo da implementação.
Testes automatizados demonstrando o funcionamento.
Utilize docker/docker-compose para que possamos realizar os testes de sua aplicação.
Deploy realizado no Google Cloud Run (free tier) e endereço ativo para ser acessado.
```
## Deployment

Passos para o deploy:
- Criado um artifact registry no GCP
- build local da aplicação com:
```docker build . -t us-central1-docker.pkg.dev/laboratorio-curso-fc/go-expert/weater```
- push da imagem para o registry:
```docker push us-central1-docker.pkg.dev/laboratorio-curso-fc/go-expert/weater:latest```
- deploy da aplicação no cloud run:
  - a variavel WEATHERAPI_KEY é a API Key para https://www.weatherapi.com/
```
gcloud run deploy \
    go-expert-lab-cloud-run \
    --allow-unauthenticated \
    --region=us-central1 \
    --image us-central1-docker.pkg.dev/laboratorio-curso-fc/go-expert/weater:latest \
    --port 8080 \
    --set-env-vars "WEATHERAPI_KEY=$WEATHERAPI_KEY"
```
resultado:
```
➜  fc-pos-go-labs-deploy-cloudrun git:(main) gcloud run deploy \
    go-expert-lab-cloud-run \
    --allow-unauthenticated \
    --region=us-central1 \
    --image us-central1-docker.pkg.dev/laboratorio-curso-fc/go-expert/weater:latest \
    --port 8080 \
    --set-env-vars "WEATHERAPI_KEY=$WEATHERAPI_KEY"
Deploying container to Cloud Run service [go-expert-lab-cloud-run] in project [laboratorio-curso-fc] region [us-central1]
✓ Deploying new service... Done.                                                                                                                                                                                                                    
  ✓ Creating Revision...                                                                                                                                                                                                                            
  ✓ Routing traffic...                                                                                                                                                                                                                              
  ✓ Setting IAM Policy...                                                                                                                                                                                                                           
Done.                                                                                                                                                                                                                                               
Service [go-expert-lab-cloud-run] revision [go-expert-lab-cloud-run-00001-8tz] has been deployed and is serving 100 percent of traffic.
Service URL: https://go-expert-lab-cloud-run-214222014794.us-central1.run.app
```
## Testes




