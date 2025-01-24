# EmailCleaner

**EmailCleaner** é uma ferramenta desenvolvida em Go que permite processar, filtrar e limpar listas de endereços de email presentes em arquivos de texto. Ideal para administradores de sistemas, profissionais de marketing e desenvolvedores que precisam organizar e validar grandes volumes de emails de forma eficiente.

## Funcionalidades

- **Filtragem Avançada**: Remove emails de domínios indesejados e aqueles que contêm palavras-chave específicas.
- **Validação de Emails**: Utiliza expressões regulares para garantir a formatação correta dos endereços de email.
- **Remoção de Duplicatas**: Elimina emails repetidos para garantir uma lista limpa e única.
- **Processamento Concorrente**: Utiliza goroutines e sincronização para acelerar o processamento de múltiplos arquivos simultaneamente.
- **Logs Detalhados**: Fornece informações sobre o processamento de cada arquivo e eventuais erros encontrados.

## Requisitos

- [Go](https://golang.org/) 1.16 ou superior

## Instalação

1. **Clone o repositório:**

   ```bash
   git clone https://github.com/rafaelwdornelas/EmailCleaner.git
