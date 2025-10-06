<div align="center">
  <img src="frontend/static/img/logos/svg/256/nobg/black.svg" alt="sik6 logo" width="150"/>
</div>

<h3 align="center">
  <img src="https://img.shields.io/badge/Go-1.25-00ADD8?logo=go&logoColor=white" alt="Go Badge" />
  <img src="https://img.shields.io/badge/Angular-20-DD0031?logo=angular&logoColor=white" alt="Angular Badge" />
  <img src="https://img.shields.io/badge/PostgreSQL-18-336791?logo=postgresql&logoColor=white" alt="Postgres Badge" />
</h3>

<p align="center">
    <strong>🍴 sik6</strong> is an application designed to simplify everyday meal preparation.<br/> 
    It suggests quick and tailored ideas to help those who often lack inspiration.<br/>
    The project aims to provide a functional base for personalized meal suggestions  
    and to evolve into a <em>modern community hub</em> for sharing and enriching the experience.<br/>
    <strong>In short</strong>: a tool built to save time and centralize culinary inspiration from around the world.
</p>

## 📖 Table of Contents

1. [Quick Start](#-quick-start)
2. [Contributing](#-contributing)
3. [Useful Links](#-useful-links)

## ⚡ Quick Start

### Installation

> Make sure you have the following locally 👇

- [Go 1.25](https://go.dev/dl/)
- [Angular 20](https://angular.dev/installation)
- [PostgreSQL 18](https://www.postgresql.org/download/)

> And now you can get the Repository 👇

```bash
git clone https://github.com/0xm0-v1/sik6
```

### Alternative with Docker 

1. Install Docker with the Compose plugin (24.x or newer).
2. Copy `.env.example` to `.env.development` and fill in the values. Set a strong `API_TOKEN` (clients must send `Authorization: Bearer <token>` for POST/PATCH/DELETE requests).
3. From the repository root:
   ```bash
   cd backend
   make up
   ```
4. Wait for the healthchecks (`make ps` shows the three services as `healthy`).
5. Verify:
   ```bash
   curl -i http://localhost:8080/livez
   curl -i http://localhost:8080/readyz
   # Angular dev server
   open http://localhost:4200/
   ```
6. Useful commands:
   - `make down` — stop and remove the containers
   - `make logs` — follow the stack logs
   - `make up-detach` — run the stack in the background

> The containers mount the local source tree (`backend/`, `frontend/`) so code changes are picked up immediately (Go is reloaded by Air, Angular runs in watch mode).

## 🔗 Useful Links

- 🎨 [Designs (Figma)](https://www.figma.com/files/team/1304730716403620532/project/127010998/sik6?fuid=1254181321735275578) – UI/UX mockups and design system

## 🤝 Contributing

_Please read the [CONTRIBUTING.md](./CONTRIBUTING.md) before opening a pull request._
