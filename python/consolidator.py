import redis

# Conexión a Redis en el namespace kafka
r = redis.Redis(host='redis-service.kafka.svc.cluster.local', port=6379, decode_responses=True)

# Buscar todas las claves que coincidan
keys = r.keys('kafka:my.topic:offset:*')

# Limpiar el hash consolidado
r.delete('kafka:my.topic:data')

# Consolidar claves en un único hash
for key in keys:
    value = r.get(key)
    if value:
        offset = key.split(':')[-1]
        r.hset('kafka:my.topic:data', offset, value)

print("✅ Consolidación completada.")
