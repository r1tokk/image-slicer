#!/bin/bash

FREEZE_DIR="./freeze_output"
SLICED_DIR="./sliced_output"

EXCLUDE_LIST=(
    ".git"
    "node_modules"
    "vendor"
    ".env"
    ".idea"
    "freeze_output"
    "sliced_output"
    "image-slicer"
    "\.png$"
    "\.jpg$"
)

if ! command -v freeze &> /dev/null; then
    echo "Помилка: утиліту 'freeze' не знайдено. Встанови її перед запуском."
    exit 1
fi

if ! command -v image-slicer &> /dev/null; then
    echo "Помилка: утиліту 'image-slicer' не знайдено у PATH."
    echo "Переконайся, що вона скомпільована (go build) та доступна для виклику."
    exit 1
fi

mkdir -p "$FREEZE_DIR"
mkdir -p "$SLICED_DIR"

EXCLUDE_PATTERN=$(IFS="|"; echo "${EXCLUDE_LIST[*]}")

echo "Починаємо масову генерацію та нарізку..."

find . -type f -print0 | while IFS= read -r -d '' filepath; do
    
    clean_path="${filepath#./}"

    if echo "$clean_path" | grep -qE "($EXCLUDE_PATTERN)"; then
        continue
    fi

    safe_name=$(echo "$clean_path" | tr '/' '_')
    freeze_file="$FREEZE_DIR/freeze_${safe_name}.png"

    echo "Freeze: $clean_path -> $freeze_file"
    
    if freeze "$filepath" -o "$freeze_file" > /dev/null 2>&1; then
        echo "🔪 Slicing: $freeze_file"
        
        image-slicer "$freeze_file" --output "$SLICED_DIR"
        
    else
        echo "Помилка генерації freeze для $clean_path. Пропускаємо."
    fi

done

echo "Готово! "
echo "Цілі скріншоти: $FREEZE_DIR"
echo "Нарізані фрагменти: $SLICED_DIR"