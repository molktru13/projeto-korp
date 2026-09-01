# Desafio DevOps - Projeto Korp

Este repositório contém a solução do Desafio DevOps da Korp.

## Estrutura do Projeto

- `app/`: Código fonte em Go e Dockerfile.
- `docker/`: Configurações do Docker Compose, Nginx e Prometheus.
- `grafana/`: Configurações de provisionamento (dashboards e datasources) automatizadas para o Grafana.
- `ansible/`: Playbook Ansible que automatiza toda a instalação, build, configuração e validação.

## Como Executar

Tudo está automatizado pelo Ansible. Para provisionar o ambiente:

```bash
cd ansible
ansible-playbook -i inventory playbook.yml -K
```

Após a execução, o serviço estará disponível em `http://localhost:80/projeto-korp`.
O dashboard do Grafana estará em `http://localhost:3000` (usuário e senha configurados na primeira execução dependem da variável ambiente, mas por padrão admin/admin, o dashboard `Dashboard Projeto Korp` é provisionado automaticamente).

## Features
- **Multi-stage Dockerfile**: Imagem enxuta (`scratch`), otimizando segurança e tamanho.
- **Observabilidade**: Métricas expostas via prometheus client, scrape do Prometheus e Dashboards criados automaticamente no Grafana.
- **Automação 100%**: Único comando para instalar dependências do Docker (se não existirem), construir imagens, subir serviços e validar o funcionamento.
