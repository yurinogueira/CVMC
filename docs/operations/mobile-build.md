---
type: runbook
title: "Empacotamento Mobile Android e Guia de Build com Capacitor — CVMC"
description: "Guia canônico e operacional para compilação, sincronização, geração de APK de debug e Android App Bundle (.aab) para a Google Play Store no CVMC."
tags:
  - mobile
  - android
  - capacitor
  - build
  - google-play
  - runbook
timestamp: 2026-10-11
---

# 📱 Empacotamento Mobile Android e Guia de Build com Capacitor

Este documento descreve o fluxo canônico de engenharia e operações para desenvolvimento, sincronização e compilação do aplicativo móvel nativo **CVMC (Como Vai Meu Carro)** para a plataforma **Android**, utilizando o ecossistema **Capacitor 8**.

---

## 🏗️ Arquitetura do Shell Nativo Mobile

O aplicativo Android do CVMC opera através de um shell nativo leve acoplado à SPA React 19 do diretório `frontend/`:

```mermaid
flowchart TD
    subgraph Frontend["SPA React 19 (/frontend)"]
        UI["Interface Material UI v6<br>Safe Area Insets & Responsive"]
        BRIDGE["MobileAppBridge<br>(StatusBar, SplashScreen, BackButton)"]
        API["apiClient (Axios)<br>Detecção Nativa + Bearer Fallback"]
        STORE["Zustand useAuthStore<br>@capacitor/preferences"]
    end

    subgraph Native["Plataforma Nativa Android (/frontend/android)"]
        WEBVIEW["Android System WebView<br>Scheme: https, cleartextTraffic: false"]
        PLUGINS["Capacitor Plugins Nativo<br>@capacitor/app<br>@capacitor/status-bar<br>@capacitor/splash-screen<br>@capacitor/preferences"]
        GRADLE["Gradle Build System<br>Debug APK / Release AAB"]
    end

    subgraph Backend["API CVMC"]
        SERVER["https://api-cvmc.yurinogueira.dev.br/api/v1<br>Bearer Token / HttpOnly Cookies"]
    end

    UI --> BRIDGE
    BRIDGE --> PLUGINS
    API --> STORE
    STORE --> PLUGINS
    API -->|HTTPS| SERVER
    WEBVIEW --> UI
```

### Configurações de Identificação e Segurança
- **Identificador do Aplicativo (App ID)**: `com.cvmc.app`
- **Nome da Aplicação**: `CVMC`
- **Diretório de Web Assets**: `dist`
- **Segurança de Rede**: `android:usesCleartextTraffic="false"` no `AndroidManifest.xml`, garantindo tráfego 100% via HTTPS.
- **Autenticação em WebView**: Detecção automática de ambiente nativo com `Capacitor.isNativePlatform()`. Os tokens de sessão (`accessToken` e `refreshToken`) são persistidos de maneira segura via `@capacitor/preferences` e transmitidos com o cabeçalho `Authorization: Bearer <token>`.

---

## 🛠️ Pré-Requisitos de Ambiente

Para gerar os pacotes nativos compilados (`.apk` ou `.aab`), a máquina ou servidor de build deve possuir:

1. **Node.js** >= 20 e **npm** >= 10.
2. **Java Development Kit (JDK)**: JDK 17 ou 21 (recomendado Azul Zulu ou OpenJDK).
3. **Android SDK / Android Studio**:
   - `ANDROID_HOME` configurado (ex: `/home/usuario/Android/Sdk` ou `C:\Users\usuario\AppData\Local\Android\Sdk`).
   - Android SDK Command-line Tools e Build-Tools (API 34 ou 35).
   - Android Platform SDK 35 (`android-35`).

---

## ⚡ Comandos e Scripts Rápidos de Desenvolvimento

Todos os scripts estão mapeados em `frontend/package.json`:

```bash
# 1. Compilar o frontend web e sincronizar com o projeto nativo Android
npm --prefix frontend run cap:build

# 2. Sincronizar plugins e cópia rápida de assets após build web
npm --prefix frontend run cap:sync

# 3. Abrir o projeto diretamente no Android Studio
npm --prefix frontend run cap:open
```

---

## 📦 Geração de Builds via Terminal (Gradle)

Navegue até a pasta nativa `frontend/android`:

### 1. Build de Debug (APK de Teste Local)
```bash
cd frontend/android
./gradlew assembleDebug
```
- **Localização do artefato gerado**:
  `frontend/android/app/build/outputs/apk/debug/app-debug.apk`

### 2. Build de Produção (Android App Bundle - AAB para Google Play)
```bash
cd frontend/android
./gradlew bundleRelease
```
- **Localização do artefato gerado**:
  `frontend/android/app/build/outputs/bundle/release/app-release.aab`

---

## 🔑 Assinatura e Keystore para Publicação na Google Play

Para publicar na Google Play Store, o pacote `.aab` deve ser assinado com uma chave privada RSA de 2048+ bits:

### 1. Geração da Chave Privada (Keystore)
> [!CAUTION]
> Guarde o arquivo `cvmc-release-key.jks` e suas senhas em um cofre seguro (ex: 1Password, Bitwarden ou Secrets do repositório). A perda desta chave impede o envio de atualizações para a Google Play Store.

```bash
keytool -genkey -v \
  -keystore cvmc-release-key.jks \
  -alias cvmc-key \
  -keyalg RSA \
  -keysize 2048 \
  -validity 10000
```

### 2. Configuração de Variáveis ou `gradle.properties`
Em ambiente de CI/CD ou local seguro, configure no `~/.gradle/gradle.properties`:

```properties
CVMC_RELEASE_STORE_FILE=/caminho/seguro/cvmc-release-key.jks
CVMC_RELEASE_KEY_ALIAS=cvmc-key
CVMC_RELEASE_STORE_PASSWORD=senha_do_keystore
CVMC_RELEASE_KEY_PASSWORD=senha_da_chave
```

---

## 🤖 Compilação Automatizada no GitHub Actions

A pipeline de Frontend (`.github/workflows/frontend.yml`) compila automaticamente o APK de Debug a cada push na branch `main` e em Pull Requests:
1. Constrói a SPA com as variáveis de produção (`VITE_GOOGLE_CLIENT_ID` e `VITE_API_BASE_URL`).
2. Sincroniza os web assets com o projeto nativo via `npx cap sync android`.
3. Executa o build Gradle `./gradlew assembleDebug`.
4. Disponibiliza o arquivo `app-debug-apk` na aba **Artifacts** da execução no GitHub Actions para download e teste direto no celular.

---

## 🔐 Autenticação Social Google no Android (Google Credential Manager)

O Google bloqueia explicitamente fluxos de autenticação web baseados em popups/postMessage dentro de WebViews embutidas no Android (`disallowed_useragent`). Para viabilizar login e vínculo social no app Android, o CVMC adota o plugin nativo `@capgo/capacitor-social-login`, acoplado ao moderno **Android Credential Manager** (`androidx.credentials`):

### 1. Configuração no Google Cloud Console
No mesmo projeto do Google Cloud onde reside o Client ID Web (`ID do cliente para Aplicativo da Web`):
1. Acesse **APIs e Serviços > Credenciais > Criar Credenciais > ID do cliente OAuth**.
2. Selecione o tipo de aplicativo **Android**.
3. Preencha o nome do pacote: `com.cvmc.app`.
4. Adicione a impressão digital do certificado SHA-1:
   - **Ambiente de Testes / Debug Local**:
     ```
     D6:0C:13:47:7D:E7:29:43:DF:AA:FB:B3:F5:10:00:BC:5A:7D:AF:02
     ```
     *(Gerado a partir do `~/.android/debug.keystore` local)*.
   - **Ambiente de Produção**:
     O SHA-1 da chave privada `cvmc-release-key.jks` ou da chave de assinatura do app gerada pelo Google Play App Signing.
5. O Web Client ID continua sendo informado na variável `VITE_GOOGLE_CLIENT_ID` (frontend) e `GOOGLE_CLIENT_ID` (backend Go), sendo repassado como `webClientId` para o Credential Manager gerar um ID Token válido para a audiência da API.

---

## 🧪 Validação e Checklist Pré-Publicação

Antes de submeter uma nova versão à Google Play:
- [ ] Executar `./scripts/check.sh all` na raiz do projeto (100% verde).
- [ ] Garantir que o `versionCode` e `versionName` foram incrementados em `frontend/android/app/build.gradle`.
- [ ] Validar a navegação via hardware back button no dispositivo físico ou emulador Android.
- [ ] Validar autenticação e silent refresh funcionando via conexão externa HTTPS.
- [ ] Validar login e vínculo com conta Google via Credential Manager nativo.
- [ ] Testar Safe Area Insets em aparelhos com entalhe/camera cutout e barra de gestos.
