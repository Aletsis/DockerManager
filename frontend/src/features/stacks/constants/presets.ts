import type { ComposePreset } from '../../../types';

export const COMPOSE_PRESETS: ComposePreset[] = [
  {
    id: 'minimal',
    name: 'Mínimo (Nginx)',
    description: 'Servidor web Nginx Alpine básico con mapeo de puerto 8080.',
    defaultName: 'nginx-web',
    yaml: `services:
  web:
    image: nginx:alpine
    container_name: web-app
    restart: unless-stopped
    ports:
      - "8080:80"
`,
  },
  {
    id: 'postgres-redis',
    name: 'PostgreSQL 16 + Redis + Adminer',
    description: 'Base de datos relacional, caché en memoria y gestor web Adminer.',
    defaultName: 'db-stack',
    yaml: `services:
  postgres:
    image: postgres:16-alpine
    container_name: postgres-db
    restart: unless-stopped
    environment:
      POSTGRES_DB: mydatabase
      POSTGRES_USER: user
      POSTGRES_PASSWORD: secretpassword
    ports:
      - "5432:5432"
    volumes:
      - pgdata:/var/lib/postgresql/data

  redis:
    image: redis:alpine
    container_name: redis-cache
    restart: unless-stopped
    ports:
      - "6379:6379"

  adminer:
    image: adminer:latest
    container_name: adminer-ui
    restart: unless-stopped
    ports:
      - "8088:8080"
    depends_on:
      - postgres

volumes:
  pgdata:
`,
  },
  {
    id: 'node-mongo',
    name: 'Node.js + MongoDB + Mongo Express',
    description: 'Stack backend Node con base de datos NoSQL y panel de administración web.',
    defaultName: 'node-mongo-stack',
    yaml: `services:
  app:
    image: node:20-alpine
    container_name: node-api
    restart: unless-stopped
    working_dir: /app
    command: node -e "require('http').createServer((req, res) => res.end('Node.js Stack Running!')).listen(3000)"
    ports:
      - "3000:3000"
    depends_on:
      - mongo

  mongo:
    image: mongo:latest
    container_name: mongodb
    restart: unless-stopped
    ports:
      - "27017:27017"
    environment:
      MONGO_INITDB_ROOT_USERNAME: root
      MONGO_INITDB_ROOT_PASSWORD: rootpassword
    volumes:
      - mongodata:/data/db

  mongo-express:
    image: mongo-express:latest
    container_name: mongo-express-ui
    restart: unless-stopped
    ports:
      - "8081:8081"
    environment:
      ME_CONFIG_MONGODB_ADMINUSERNAME: root
      ME_CONFIG_MONGODB_ADMINPASSWORD: rootpassword
      ME_CONFIG_MONGODB_SERVER: mongo
    depends_on:
      - mongo

volumes:
  mongodata:
`,
  },
  {
    id: 'wordpress-mysql',
    name: 'WordPress + MySQL 8',
    description: 'Instalación de WordPress con base de datos MySQL y almacenamiento persistente.',
    defaultName: 'wordpress-site',
    yaml: `services:
  wordpress:
    image: wordpress:latest
    container_name: wordpress-app
    restart: unless-stopped
    ports:
      - "8000:80"
    environment:
      WORDPRESS_DB_HOST: db:3306
      WORDPRESS_DB_USER: wp_user
      WORDPRESS_DB_PASSWORD: wp_password
      WORDPRESS_DB_NAME: wp_database
    volumes:
      - wp_data:/var/www/html
    depends_on:
      - db

  db:
    image: mysql:8.0
    container_name: wordpress-db
    restart: unless-stopped
    environment:
      MYSQL_DATABASE: wp_database
      MYSQL_USER: wp_user
      MYSQL_PASSWORD: wp_password
      MYSQL_RANDOM_ROOT_PASSWORD: '1'
    volumes:
      - db_data:/var/lib/mysql

volumes:
  wp_data:
  db_data:
`,
  },
];
