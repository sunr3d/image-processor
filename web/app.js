let currentImageId = null;

// Загрузка изображения
document.getElementById('upload-form').addEventListener('submit', async (e) => {
    e.preventDefault();
    
    const fileInput = document.getElementById('image-input');
    const file = fileInput.files[0];
    
    if (!file) return;
    
    const formData = new FormData();
    formData.append('image', file);
    
    try {
        const response = await fetch('/upload', {
            method: 'POST',
            body: formData
        });
        
        if (!response.ok) {
            throw new Error('Ошибка загрузки');
        }
        
        const result = await response.json();
        currentImageId = result.id;
        
        showResults();
        showImageId();
        showStatus('Изображение загружено, идет обработка...');
        checkStatus();
        
    } catch (error) {
        alert('Ошибка: ' + error.message);
    }
});

// Поиск по ID
document.getElementById('search-form').addEventListener('submit', async (e) => {
    e.preventDefault();
    
    const searchInput = document.getElementById('search-input');
    const imageId = searchInput.value.trim();
    
    if (!imageId) return;
    
    try {
        const response = await fetch(`/status/${imageId}`);
        
        if (response.ok) {
            currentImageId = imageId;
            showResults();
            showImageId();
            checkStatus();
        } else if (response.status === 404) {
            alert('Изображение с таким ID не найдено');
        } else {
            alert('Ошибка при поиске изображения');
        }
        
    } catch (error) {
        alert('Ошибка: ' + error.message);
    }
});

// Проверка статуса
async function checkStatus() {
    if (!currentImageId) return;
    
    try {
        const response = await fetch(`/status/${currentImageId}`);
        
        if (response.ok) {
            const status = await response.json();
            
            if (status.status === 'completed') {
                showStatus('Обработка завершена!');
                showImages();
            } else if (status.status === 'failed') {
                showStatus('Ошибка обработки: ' + status.message);
            } else {
                showStatus('Обработка...');
                setTimeout(checkStatus, 2000);
            }
        } else {
            showStatus('Ошибка проверки статуса');
        }
        
    } catch (error) {
        showStatus('Ошибка проверки статуса');
    }
}

// Показать изображения
async function showImages() {
    const imagesDiv = document.getElementById('images');
    imagesDiv.innerHTML = '';
    
    const types = ['original', 'resized', 'thumbnail', 'watermarked'];
    
    for (const type of types) {
        try {
            const response = await fetch(`/image/${currentImageId}?type=${type}`);
            if (response.ok) {
                const blob = await response.blob();
                const url = URL.createObjectURL(blob);
                
                const img = document.createElement('img');
                img.src = url;
                img.style.maxWidth = '300px';
                img.style.margin = '10px';
                
                const div = document.createElement('div');
                div.innerHTML = `<h3>${type}</h3>`;
                div.appendChild(img);
                
                imagesDiv.appendChild(div);
            }
        } catch (error) {
            console.log(`Ошибка загрузки ${type}:`, error);
        }
    }
    
    document.getElementById('delete-btn').style.display = 'block';
}

// Показать ID изображения
function showImageId() {
    const idDisplay = document.getElementById('image-id-display');
    idDisplay.innerHTML = `<strong>ID изображения:</strong> ${currentImageId}`;
}

// Показать статус
function showStatus(message) {
    document.getElementById('status').textContent = message;
}

// Показать результаты
function showResults() {
    document.getElementById('upload-section').style.display = 'none';
    document.getElementById('search-section').style.display = 'none';
    document.getElementById('results-section').style.display = 'block';
}

// Удаление изображения
document.getElementById('delete-btn').addEventListener('click', async () => {
    if (!currentImageId) return;
    
    try {
        const response = await fetch(`/image/${currentImageId}`, {
            method: 'DELETE'
        });
        
        if (response.ok) {
            alert('Изображение удалено');
            resetForm();
        } else {
            alert('Ошибка удаления');
        }
        
    } catch (error) {
        alert('Ошибка: ' + error.message);
    }
});

// Кнопка "Назад"
document.getElementById('back-btn').addEventListener('click', () => {
    resetForm();
});

// Сброс формы
function resetForm() {
    currentImageId = null;
    document.getElementById('upload-section').style.display = 'block';
    document.getElementById('search-section').style.display = 'block';
    document.getElementById('results-section').style.display = 'none';
    document.getElementById('image-input').value = '';
    document.getElementById('search-input').value = '';
    document.getElementById('delete-btn').style.display = 'none';
}