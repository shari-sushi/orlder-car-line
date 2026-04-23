.PHONY: ngrok-up make-db-up

ngrok-up:
	(cd app && npm run dev) & npx ngrok http 3000

make-db-up:
  docker-compose up