#!/bin/bash

# Acesso do banco
DB_HOST="localhost"
DB_PORT="5432"
DB_USER="postgres"
DB_NAME="atmos"

# Cores para deixar mais intuitivo
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${YELLOW}🗄️  Sistema de Migrations Simples${NC}"
echo "=================================="

# Função para criar nova migration
# Essa função foi o motivo para criar esse script.
nova_migration() {
    if [ -z "$1" ]; then
        echo -e "${RED}❌ Erro: Descrição vazia!${NC}"
        echo "Exemplo: ./migrate.sh nova 'adicionar coluna histórico'"
        exit 1
    fi
    
    # Pradronização dos nomes das migrations
    NUMERO=$(ls migrations/*.sql | wc -l | tr -d ' ')
    NUMERO=$(printf "%03d" $((NUMERO)))
    DESCRICAO=$(echo "$1" | tr ' ' '_' | tr '[:upper:]' '[:lower:]')
    ARQUIVO="migrations/${NUMERO}_${DESCRICAO}.sql"
    
    # Registro da migration e código SQL
    cat > "$ARQUIVO"
    << EOF
-- Migration ${NUMERO}: $1
-- Data: $(date +%d/%m/%Y)

-- Query SQL

EOF
    
    echo -e "${GREEN}✅ Migration criada: $ARQUIVO${NC}"
    echo -e "${YELLOW}💡 Próximos passos:${NC}"
    echo "   1. Edite o arquivo $ARQUIVO"
    echo "   2. Execute: ./migrate.sh aplicar após editar."
}

# Função para aplicar migrations pendentes
# Basicamente roda pelos arquivos ".sql" exceto o 000_setup, verifica se já foi executado, caso negativo aplica a migration e caso positivo informa que já foi aplicada.
aplicar_migrations() {
    echo -e "${YELLOW}🔄 Verificando migrations pendentes...${NC}"
    
    for arquivo in migrations/*.sql; do
        if [ -f "$arquivo" ]; then
            nome_arquivo=$(basename "$arquivo")
            
            if [[ "$nome_arquivo" == "000_setup_migrations.sql" ]]; then
                continue
            fi
            
            resultado=$(psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "SELECT COUNT(*) FROM schema_migrations WHERE migration_name = '$nome_arquivo';" | tr -d ' ')
            
            if [ "$resultado" = "0" ]; then
                echo -e "${YELLOW}⏳ Aplicando: $nome_arquivo${NC}"
                
                psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f "$arquivo"
                
                if [ $? -eq 0 ]; then
                    psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "INSERT INTO schema_migrations (migration_name) VALUES ('$nome_arquivo');"
                    echo -e "${GREEN}✅ $nome_arquivo aplicada com sucesso!${NC}"
                else
                    echo -e "${RED}❌ Erro ao aplicar $nome_arquivo${NC}"
                    exit 1
                fi
            else
                echo -e "${GREEN}✓${NC} $nome_arquivo (já aplicada)"
            fi
        fi
    done
    
    echo -e "${GREEN}🎉 Todas as migrations estão atualizadas!${NC}"
}

# Função para ver o estado
# Depois do primeiro erro precisei fazer alguma forma de verificar o estado que ficou após o erro, esse foi o resultado
status() {
    echo -e "${YELLOW}📊 Status das Migrations:${NC}"
    echo "========================"
    
    for arquivo in migrations/*.sql; do
        if [ -f "$arquivo" ]; then
            nome_arquivo=$(basename "$arquivo")
            
            if [[ "$nome_arquivo" == "000_setup_migrations.sql" ]]; then
                continue
            fi
            
            resultado=$(psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "SELECT executed_at FROM schema_migrations WHERE migration_name = '$nome_arquivo';" | tr -d ' ')
            
            if [ -n "$resultado" ] && [ "$resultado" != "" ]; then
                echo -e "${GREEN}✅${NC} $nome_arquivo (executada em: $resultado)"
            else
                echo -e "${RED}❌${NC} $nome_arquivo (pendente)"
            fi
        fi
    done
}

# Menu Visual
# a principio decidi inserir-lo no próprio script, mas talvez eu coloque no README.md
case "$1" in
    "nova")
        nova_migration "$2"
        ;;
    "aplicar")
        aplicar_migrations
        ;;
    "status")
        status
        ;;
    *)
        echo -e "${YELLOW}Como usar:${NC}"
        echo "  ./migrate.sh nova 'descrição'    - Criar nova migration"
        echo "  ./migrate.sh aplicar             - Aplicar migrations pendentes"
        echo "  ./migrate.sh status              - Ver status das migrations"
        echo ""
        echo -e "${YELLOW}Exemplo:${NC}"
        echo "  ./migrate.sh nova 'adicionar tabela histórico'"
        echo "  ./migrate.sh aplicar"
        ;;
esac
EOF