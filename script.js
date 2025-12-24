let chartInstance = null;

async function startBenchmark() {
    const nInput = document.getElementById('inputN');
    const xInput = document.getElementById('inputX');
    
    const n = parseInt(nInput.value);
    const x = parseInt(xInput.value);

    // Validasi input
    if (!n || x < 0 || x > n) {
        alert("Harap masukkan angka yang valid (X harus <= N)");
        return;
    }

    // Tampilkan Loading
    document.getElementById('results').classList.add('hidden');
    document.getElementById('loading').classList.remove('hidden');

    try {
        // Panggil API Go
        const response = await fetch('/api/benchmark', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ n: n, x: x })
        });

        if (!response.ok) throw new Error("Gagal mengambil data dari server");

        const data = await response.json();

        // Tampilkan Data
        renderStats(data);
        renderChart(data);
        renderSamples(data.samples);

        document.getElementById('loading').classList.add('hidden');
        document.getElementById('results').classList.remove('hidden');

    } catch (error) {
        console.error("Error:", error);
        alert("Terjadi kesalahan. Pastikan server Go (main.go) sudah berjalan.");
        document.getElementById('loading').classList.add('hidden');
    }
}

function renderStats(data) {
    document.getElementById('iterTime').innerText = data.iterativeTime;
    document.getElementById('recTime').innerText = data.recursiveTime;
}

function renderChart(data) {
    const ctx = document.getElementById('benchmarkChart').getContext('2d');

    if (chartInstance) {
        chartInstance.destroy();
    }

    chartInstance = new Chart(ctx, {
        type: 'line',
        data: {
            labels: data.graphLabels,
            datasets: [
                {
                    label: 'Iteratif',
                    data: data.graphIterative,
                    borderColor: '#5c85f7',
                    backgroundColor: '#5c85f7',
                    tension: 0.3
                },
                {
                    label: 'Rekursif',
                    data: data.graphRecursive,
                    borderColor: '#f7a05c',
                    backgroundColor: '#f7a05c',
                    tension: 0.3
                }
            ]
        },
        options: {
            responsive: true,
            plugins: {
                legend: { labels: { color: '#ccc' } }
            },
            scales: {
                x: { 
                    ticks: { color: '#888' }, 
                    grid: { color: '#333' },
                    title: { display: true, text: 'Jumlah Data (N)', color:'#666'}
                },
                y: { 
                    ticks: { color: '#888' }, 
                    grid: { color: '#333' },
                    title: { display: true, text: 'Waktu (ms)', color:'#666'}
                }
            }
        }
    });
}

function renderSamples(samples) {
    const grid = document.getElementById('plateGrid');
    grid.innerHTML = ""; 

    if (!samples || samples.length === 0) {
        grid.innerHTML = "<p style='color:#555'>Tidak ada palindrom ditemukan.</p>";
        return;
    }

    samples.forEach((plate, index) => {
        const div = document.createElement('div');
        div.className = 'plate-item';
        div.innerHTML = `
            <div>${plate}</div>
            <div style="font-size:0.7rem; color:#555; margin-top:4px">#${index+1}</div>
        `;
        grid.appendChild(div);
    });
}