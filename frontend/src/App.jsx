import React, { useState, useEffect } from 'react';

const API_BASE = 'http://localhost:8080';

function App() {
  const [grid, setGrid] = useState([]);
  const [width, setWidth] = useState(21);
  const [height, setHeight] = useState(21);
  const [path, setPath] = useState([]);
  const [isSolving, setIsSolving] = useState(false);
  const [difficulty, setDifficulty] = useState('Medium');
  const [startNode, setStartNode] = useState({ x: 0, y: 0 });
  const [endNode, setEndNode] = useState({ x: 20, y: 20 });
  const [isMouseDown, setIsMouseDown] = useState(false);
  const [paintMode, setPaintMode] = useState(1); // 1 = wall, 0 = path
  const [notification, setNotification] = useState(null);
  const [notificationType, setNotificationType] = useState('success');

  useEffect(() => {
    generateInitialGrid(21, 21);

    const handleGlobalMouseUp = () => setIsMouseDown(false);
    window.addEventListener('mouseup', handleGlobalMouseUp);
    return () => window.removeEventListener('mouseup', handleGlobalMouseUp);
  }, []);

  const showToast = (msg, type = 'success') => {
    setNotification(msg);
    setNotificationType(type);
    setTimeout(() => {
      setNotification(null);
      setNotificationType('success');
    }, 3000);
  };

  const generateInitialGrid = (w, h) => {
    const newGrid = Array(h).fill().map(() => Array(w).fill(0));
    setGrid(newGrid);
    setPath([]);
    setStartNode({ x: 0, y: 0 });
    setEndNode({ x: w - 1, y: h - 1 });
  };

  const handleDifficultyChange = (e) => {
    const val = e.target.value;
    setDifficulty(val);
    let w = 21, h = 21;

    if (val === 'Easy') {
      w = 11; h = 11;
    } else if (val === 'Medium') {
      w = 21; h = 21;
    } else if (val === 'Hard') {
      w = 51; h = 21;
    }

    setWidth(w);
    setHeight(h);
    generateInitialGrid(w, h);
  };

  const generateRandomMaze = async () => {
    try {
      const resp = await fetch(`${API_BASE}/generate`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ width, height, difficulty })
      });
      const data = await resp.json();

      // Preserve current start/end if they are within new bounds
      let newStart = { ...startNode };
      let newEnd = { ...endNode };

      // Safety check: ensure coordinates fit in new dimensions
      if (newStart.x >= width) newStart.x = width - 1;
      if (newStart.y >= height) newStart.y = height - 1;
      if (newEnd.x >= width) newEnd.x = width - 1;
      if (newEnd.y >= height) newEnd.y = height - 1;

      // Force Start and End to be paths (0) in the new grid
      const newGrid = data.grid;

      const ensureConnectivity = (node, g) => {
        if (!g[node.y] || g[node.y][node.x] === undefined) return;
        g[node.y][node.x] = 0; // Force self to path

        // Check if already connected to any path
        const dirs = [[0, 1], [0, -1], [1, 0], [-1, 0]];
        const isConnected = dirs.some(([dx, dy]) => {
          const ny = node.y + dy;
          const nx = node.x + dx;
          return g[ny] && g[ny][nx] === 0;
        });

        // If isolated, open just ONE neighbor to connect
        if (!isConnected) {
          // Try to find a valid neighbor to open
          for (const [dx, dy] of dirs) {
            const ny = node.y + dy;
            const nx = node.x + dx;
            if (g[ny] && g[ny][nx] !== undefined) {
              g[ny][nx] = 0;
              break; // Only open one path, then stop
            }
          }
        }
      };

      ensureConnectivity(newStart, newGrid);
      ensureConnectivity(newEnd, newGrid);

      setGrid(newGrid);
      setPath([]);
      setStartNode(newStart);
      setEndNode(newEnd);
      setEndNode(newEnd);
    } catch (err) {
      showToast('Backend Go error.', 'error');
    }
  };

  const solveMaze = async () => {
    if (isSolving) return;
    setIsSolving(true);
    setPath([]);
    try {
      const resp = await fetch(`${API_BASE}/solve`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ grid, start: startNode, end: endNode })
      });
      const data = await resp.json();
      if (data.path) animatePath(data.path);
      else {
        showToast('Tidak ada jalan!', 'error');
        setIsSolving(false);
      }
    } catch (err) {
      showToast('Gagal solve.', 'error');
      setIsSolving(false);
    }
  };

  const animatePath = (pathNodes) => {
    const ordered = [...pathNodes].reverse();
    let i = 0;
    const timer = setInterval(() => {
      if (i >= ordered.length) {
        clearInterval(timer);
        setIsSolving(false);
        showToast("Yay! Jalan keluar ditemukan!");
        return;
      }
      setPath(prev => [...prev, ordered[i]]);
      i++;
    }, 25);
  };

  const handleMouseDown = (y, x, e) => {
    if (isSolving) return;
    if (e.shiftKey || e.altKey) {
      toggleCell(y, x, e);
      return;
    }
    const newVal = grid[y][x] === 1 ? 0 : 1;
    setPaintMode(newVal);
    setIsMouseDown(true);
    updateCell(y, x, newVal);
  };

  const handleMouseEnter = (y, x) => {
    if (!isMouseDown || isSolving) return;
    updateCell(y, x, paintMode);
  };

  const updateCell = (y, x, val) => {
    if ((y === startNode.y && x === startNode.x) || (y === endNode.y && x === endNode.x)) return;
    const ng = [...grid];
    ng[y][x] = val;
    setGrid(ng);
    setPath([]);
    setStartNode({ x: 0, y: 0 });
    setEndNode({ x: width - 1, y: height - 1 });
  };

  const toggleCell = (y, x, e) => {
    if (isSolving) return;
    if (e.shiftKey) {
      setStartNode({ x, y });
      if (grid[y][x] === 1) { const ng = [...grid]; ng[y][x] = 0; setGrid(ng); }
      return;
    }
    if (e.altKey) {
      setEndNode({ x, y });
      if (grid[y][x] === 1) { const ng = [...grid]; ng[y][x] = 0; setGrid(ng); }
      return;
    }
  };

  return (
    <div className="app-container">
      <header className="header">
        <h1>Maze Solver AI</h1>
        <p>Tugas Besar UAS - Golang</p>
      </header>

      <div className="controls">
        <select value={difficulty} onChange={handleDifficultyChange}>
          <option value="Easy">Mudah (11x11)</option>
          <option value="Medium">Sedang (21x21)</option>
          <option value="Hard">Sulit (51x21)</option>
        </select>

        <button onClick={generateRandomMaze}>Generate Random</button>
        <button onClick={solveMaze} disabled={isSolving}>
          {isSolving ? 'Solving...' : 'Solve Maze'}
        </button>
        <button onClick={() => generateInitialGrid(width, height)} style={{ backgroundColor: '#475569' }}>Reset All</button>
      </div>

      <div className="maze-container">
        <div
          key={`${width}-${height}`}
          className="maze-grid"
          style={{
            width: 'min(90vw, 700px)',
            maxWidth: '100%',
            display: 'grid',
            justifyContent: 'center',
            margin: '0 auto',
            gridTemplateColumns: `repeat(${width}, 1fr)`,
            gap: width > 40 ? '1px' : '2px',
            padding: '4px',
            backgroundColor: '#334155',
            borderRadius: '8px'
          }}
        >
          {grid.map((row, y) => row.map((cell, x) => (
            <div
              key={`${y}-${x}`}
              onMouseDown={e => {
                e.preventDefault();
                handleMouseDown(y, x, e);
              }}
              onMouseEnter={() => handleMouseEnter(y, x)}
              className={`cell ${cell === 1 ? 'wall' : 'path'} ${path.some(p => p && p.x === x && p.y === y) ? 'solve' : ''} ${y === startNode.y && x === startNode.x ? 'start' : ''} ${y === endNode.y && x === endNode.x ? 'finish' : ''}`}
            />
          )))}
        </div>
      </div>

      <div className="stats">
        {path.length > 0 && `Jalur: ${path.length} langkah`}
        <p>Tip: <b>Klik/Seret</b> untuk buat dinding, <b>Shift+Klik</b>: Start (S), <b>Alt+Klik</b>: Finish (F).</p>
      </div>

      {notification && (
        <div style={{
          position: 'fixed',
          bottom: '20px',
          left: '50%',
          transform: 'translateX(-50%)',
          backgroundColor: notificationType === 'error' ? '#ef4444' : '#10b981',
          color: 'white',
          padding: '12px 24px',
          borderRadius: '8px',
          boxShadow: '0 4px 6px rgba(0,0,0,0.1)',
          zIndex: 1000,
          fontWeight: '600',
          animation: 'fadeIn 0.3s ease-out'
        }}>
          {notification}
        </div>
      )}
    </div>
  );
}

export default App;
