docker compose --profile postgres up -d
docker compose --profile mysql up -d
docker compose --profile postgres --profile mysql down
<!-- untuk menghapus penyimpanan lama sehingga database terinisialisasi ulang dari awal menggunakan password baru dari .env -->
docker compose --profile mysql down -v
docker compose --profile mysql up -d --force-recreate
<!-- untuk menghapus penyimpanan lama sehingga database terinisialisasi ulang dari awal menggunakan password baru dari .env -->
docker compose --profile postgres down -v
docker compose --profile postgres up -d --force-recreate