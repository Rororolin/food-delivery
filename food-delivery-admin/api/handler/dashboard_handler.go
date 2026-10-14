package handler

import (
	"net/http"
)

// DashboardHandler — отдаёт HTML-страницу админ-панели
type DashboardHandler struct{}

func NewDashboardHandler() *DashboardHandler {
	return &DashboardHandler{}
}

func (h *DashboardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(dashboardHTML))
}

const dashboardHTML = `<!DOCTYPE html>
<html lang="ru">
<head>
<meta charset="UTF-8"/>
<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
<title>Админ — Линуксоед</title>
<link href="https://fonts.googleapis.com/css2?family=Playfair+Display:wght@700;900&family=DM+Sans:wght@300;400;500&display=swap" rel="stylesheet"/>
<style>
  :root {
    --bg: #0f0e0c;
    --surface: #1a1814;
    --card: #211f1b;
    --accent: #e8c547;
    --accent2: #e87d47;
    --red: #e84747;
    --green: #47e8a0;
    --blue: #47b4e8;
    --text: #f0ece3;
    --muted: #7a7568;
    --border: #2e2b24;
  }

  * { margin: 0; padding: 0; box-sizing: border-box; }

  body {
    background: var(--bg);
    color: var(--text);
    font-family: 'DM Sans', sans-serif;
    min-height: 100vh;
  }

  header {
    padding: 1.2rem 2.5rem;
    display: flex; align-items: center; justify-content: space-between;
    border-bottom: 1px solid var(--border);
    position: sticky; top: 0; background: var(--bg); z-index: 100;
  }

  .logo {
    font-family: 'Playfair Display', serif;
    font-size: 1.4rem; font-weight: 900; color: var(--accent);
  }
  .logo span {
    font-family: 'DM Sans', sans-serif;
    font-size: 0.75rem; font-weight: 400;
    color: var(--muted); text-transform: uppercase;
    letter-spacing: 0.1em; margin-left: 0.6rem;
  }

  nav {
    display: flex; gap: 0;
    border-bottom: 1px solid var(--border);
    padding: 0 2.5rem;
    background: var(--surface);
  }
  nav button {
    background: none; border: none; border-bottom: 2px solid transparent;
    color: var(--muted); padding: 1rem 1.5rem;
    font-family: 'DM Sans', sans-serif; font-size: 0.9rem;
    cursor: pointer; transition: all 0.2s; letter-spacing: 0.03em;
  }
  nav button:hover { color: var(--text); }
  nav button.active { color: var(--accent); border-bottom-color: var(--accent); }

  main { padding: 2rem 2.5rem; max-width: 1300px; }

  .section { display: none; }
  .section.active { display: block; }

  .section-header {
    display: flex; align-items: center; justify-content: space-between;
    margin-bottom: 1.5rem;
  }
  .section-header h2 {
    font-family: 'Playfair Display', serif;
    font-size: 1.8rem; font-weight: 700;
  }

  /* Фильтры */
  .filters {
    display: flex; gap: 0; margin-bottom: 1.5rem; flex-wrap: wrap;
  }
  .filters select {
    padding: 0.7rem 1rem;
    background: var(--card); border: 1px solid var(--border); border-right: none;
    color: var(--text); font-family: 'DM Sans', sans-serif; font-size: 0.9rem;
    outline: none; cursor: pointer;
  }
  .filters select option { background: var(--card); }
  .filters button {
    padding: 0.7rem 1.2rem;
    background: var(--accent); color: var(--bg); border: none;
    font-family: 'DM Sans', sans-serif; font-weight: 500; font-size: 0.85rem;
    cursor: pointer; transition: background 0.2s; white-space: nowrap;
  }
  .filters button:hover { background: var(--accent2); }

  /* Таблица */
  .table-wrap {
    background: var(--card);
    border: 1px solid var(--border);
    overflow-x: auto;
  }
  table { width: 100%; border-collapse: collapse; }
  thead { background: var(--surface); }
  th {
    padding: 0.9rem 1.2rem; text-align: left;
    font-size: 0.72rem; font-weight: 500;
    text-transform: uppercase; letter-spacing: 0.08em;
    color: var(--muted); white-space: nowrap;
  }
  td {
    padding: 1rem 1.2rem; font-size: 0.88rem;
    border-top: 1px solid var(--border);
  }
  tr:hover td { background: #27241e; }

  code {
    font-family: monospace; font-size: 0.8rem;
    color: var(--muted); background: var(--surface);
    padding: 2px 6px;
  }

  /* Бейджи */
  .badge {
    display: inline-block; padding: 3px 10px;
    font-size: 0.72rem; font-weight: 600;
    text-transform: uppercase; letter-spacing: 0.05em;
  }
  .badge-pending    { background: rgba(232,197,71,0.15);  color: var(--accent); }
  .badge-assigned   { background: rgba(71,180,232,0.15);  color: var(--blue); }
  .badge-delivering { background: rgba(71,232,160,0.15);  color: var(--green); }
  .badge-delivered  { background: rgba(122,117,104,0.15); color: var(--muted); }
  .badge-cancelled  { background: rgba(232,71,71,0.15);   color: var(--red); }
  .badge-blocked    { background: rgba(232,71,71,0.15);   color: var(--red); }
  .badge-active     { background: rgba(71,232,160,0.15);  color: var(--green); }
  .badge-admin      { background: rgba(232,197,71,0.15);  color: var(--accent); }
  .badge-courier    { background: rgba(71,180,232,0.15);  color: var(--blue); }
  .badge-customer   { background: rgba(122,117,104,0.15); color: var(--muted); }

  /* Кнопки действий */
  .actions { display: flex; gap: 4px; flex-wrap: wrap; }
  .btn {
    padding: 4px 12px; border: none;
    font-family: 'DM Sans', sans-serif; font-size: 0.78rem; font-weight: 500;
    cursor: pointer; transition: opacity 0.15s; white-space: nowrap;
  }
  .btn:hover { opacity: 0.8; }
  .btn-danger    { background: var(--red);    color: #fff; }
  .btn-primary   { background: var(--blue);   color: #0a1a2a; }
  .btn-warning   { background: var(--accent2); color: #fff; }
  .btn-success   { background: var(--green);  color: #0a2a1a; }
  .btn-secondary { background: var(--border); color: var(--muted); }

  /* Тост */
  .toast {
    position: fixed; bottom: 2rem; right: 2rem;
    background: var(--card); border-left: 3px solid var(--accent);
    color: var(--text); padding: 0.9rem 1.4rem;
    font-size: 0.88rem; opacity: 0; transition: opacity 0.3s;
    z-index: 999; max-width: 340px;
    border: 1px solid var(--border); border-left: 3px solid var(--accent);
  }
  .toast.show { opacity: 1; }
  .toast.error { border-left-color: var(--red); }

  .loading {
    padding: 3rem; text-align: center;
    color: var(--muted); font-style: italic;
  }

  /* Модалки */
  .modal-overlay {
    display: none; position: fixed; inset: 0;
    background: rgba(0,0,0,0.7);
    z-index: 200; align-items: center; justify-content: center;
  }
  .modal-overlay.open { display: flex; }
  .modal {
    background: var(--card); border: 1px solid var(--border);
    padding: 2rem; width: 420px;
  }
  .modal h3 {
    font-family: 'Playfair Display', serif;
    font-size: 1.3rem; font-weight: 700;
    margin-bottom: 1.5rem; color: var(--accent);
  }
  .modal label {
    display: block; font-size: 0.72rem; text-transform: uppercase;
    letter-spacing: 0.08em; color: var(--muted); margin-bottom: 0.4rem;
  }
  .modal select, .modal input {
    width: 100%; padding: 0.8rem 1rem;
    background: var(--surface); border: 1px solid var(--border);
    color: var(--text); font-family: 'DM Sans', sans-serif;
    font-size: 0.9rem; outline: none; margin-bottom: 1.2rem;
  }
  .modal select option { background: var(--surface); }
  .modal select:focus, .modal input:focus { border-color: var(--accent); }
  .modal-actions { display: flex; gap: 8px; justify-content: flex-end; margin-top: 0.5rem; }
</style>
</head>
<body>

<header>
  <div class="logo">Линуксоед <span>Администратор</span></div>
  <div style="color:var(--muted);font-size:0.8rem">Панель управления</div>
</header>

<nav>
  <button class="active" onclick="showTab('orders', this)">Заказы</button>
  <button onclick="showTab('users', this)">Пользователи</button>
  <button onclick="showTab('couriers', this)">Курьеры</button>
</nav>

<main>
  <!-- ЗАКАЗЫ -->
  <div id="orders" class="section active">
    <div class="section-header">
      <h2>Заказы</h2>
    </div>
    <div class="filters">
      <select id="statusFilter">
        <option value="">Все статусы</option>
        <option value="pending">Ожидает</option>
        <option value="assigned">Назначен курьер</option>
        <option value="delivering">В доставке</option>
        <option value="delivered">Доставлен</option>
        <option value="cancelled">Отменён</option>
      </select>
      <button onclick="loadOrders()">Обновить</button>
    </div>
    <div class="table-wrap">
      <div id="ordersTable"><div class="loading">Загрузка...</div></div>
    </div>
  </div>

  <!-- ПОЛЬЗОВАТЕЛИ -->
  <div id="users" class="section">
    <div class="section-header">
      <h2>Пользователи</h2>
    </div>
    <div class="filters">
      <button onclick="loadUsers()">Обновить</button>
    </div>
    <div class="table-wrap">
      <div id="usersTable"><div class="loading">Загрузка...</div></div>
    </div>
  </div>

  <!-- КУРЬЕРЫ -->
  <div id="couriers" class="section">
    <div class="section-header">
      <h2>Курьеры</h2>
    </div>
    <div class="filters">
      <button onclick="loadCouriers()">Обновить</button>
    </div>
    <div class="table-wrap">
      <div id="couriersTable"><div class="loading">Загрузка...</div></div>
    </div>
  </div>
</main>

<!-- Модалка: изменить статус -->
<div class="modal-overlay" id="statusModal">
  <div class="modal">
    <h3>Изменить статус</h3>
    <label>Новый статус</label>
    <select id="newStatus">
      <option value="assigned">Assigned</option>
      <option value="delivering">Delivering</option>
      <option value="delivered">Delivered</option>
      <option value="cancelled">Cancelled</option>
    </select>
    <div class="modal-actions">
      <button class="btn btn-secondary" onclick="closeModal('statusModal')">Отмена</button>
      <button class="btn btn-primary" onclick="confirmStatusChange()">Сохранить</button>
    </div>
  </div>
</div>

<!-- Модалка: переназначить курьера -->
<div class="modal-overlay" id="courierModal">
  <div class="modal">
    <h3>Назначить курьера</h3>
    <label>Выберите курьера</label>
    <select id="newCourierID">
      <option value="">— загрузка... —</option>
    </select>
    <div class="modal-actions">
      <button class="btn btn-secondary" onclick="closeModal('courierModal')">Отмена</button>
      <button class="btn btn-primary" onclick="confirmCourierChange()">Назначить</button>
    </div>
  </div>
</div>

<div class="toast" id="toast"></div>

<script>
const API = '';
const COURIER_API = 'http://localhost:8082';
const CLIENT_API  = 'http://localhost:8080';
let currentOrderID = null;

// Русские названия статусов
const STATUS_RU = {
  'pending':    'Ожидает',
  'assigned':   'Назначен курьер',
  'delivering': 'В доставке',
  'delivered':  'Доставлен',
  'cancelled':  'Отменён',
};

function showTab(name, btn) {
  document.querySelectorAll('.section').forEach(s => s.classList.remove('active'));
  document.querySelectorAll('nav button').forEach(b => b.classList.remove('active'));
  document.getElementById(name).classList.add('active');
  btn.classList.add('active');
  if (name === 'orders')   loadOrders();
  if (name === 'users')    loadUsers();
  if (name === 'couriers') loadCouriers();
}

function statusBadge(s) {
  const label = STATUS_RU[s] || s;
  return '<span class="badge badge-' + s + '">' + label + '</span>';
}

// Таймаут для внешних запросов
async function fetchT(url, opts = {}, ms = 3000) {
  const ctrl = new AbortController();
  const t = setTimeout(() => ctrl.abort(), ms);
  try { return await fetch(url, { ...opts, signal: ctrl.signal }); }
  finally { clearTimeout(t); }
}

// ─── ЗАКАЗЫ ───────────────────────────────────────────────────────────────

async function loadOrders() {
  const status = document.getElementById('statusFilter').value;
  const url = API + '/admin/orders' + (status ? '?status=' + status : '');
  document.getElementById('ordersTable').innerHTML = '<div class="loading">Загрузка...</div>';
  try {
    const res = await fetch(url);
    const data = await res.json();
    const orders = data.orders || [];
    if (!orders.length) {
      document.getElementById('ordersTable').innerHTML = '<div class="loading">Заказов не найдено</div>';
      return;
    }
    let html = '<table><thead><tr><th>ID</th><th>Клиент</th><th>Курьер</th><th>Сумма</th><th>Оплачен</th><th>Статус</th><th>Действия</th></tr></thead><tbody>';
    for (const o of orders) {
      html += '<tr>';
      html += '<td><code>' + o.id.slice(0,8) + '…</code></td>';
      html += '<td>' + o.customer_id + '</td>';
      html += '<td>' + (o.courier_id || '—') + '</td>';
      html += '<td style="color:var(--accent)">' + o.totalAmount.toFixed(0) + ' ₽</td>';
      html += '<td>' + (o.is_paid ? '<span style="color:var(--green)">✓</span>' : '<span style="color:var(--muted)">—</span>') + '</td>';
      html += '<td>' + statusBadge(o.status) + '</td>';
      html += '<td><div class="actions">';
      if (o.status !== 'delivered' && o.status !== 'cancelled') {
        html += '<button class="btn btn-primary" onclick="openStatusModal(\'' + o.id + '\')">Статус</button>';
        html += '<button class="btn btn-warning" onclick="openCourierModal(\'' + o.id + '\')">Курьер</button>';
        html += '<button class="btn btn-danger" onclick="cancelOrder(\'' + o.id + '\')">Отменить</button>';
      }
      if (o.is_paid && !o.is_refunded && o.status === 'cancelled') {
        html += '<button class="btn btn-success" onclick="doRefund(\'' + o.id + '\')">Возврат</button>';
      }
      if (o.status === 'delivered') {
        html += '<span style="color:var(--muted);font-size:0.8rem">завершён</span>';
      }
      html += '</div></td></tr>';
    }
    html += '</tbody></table>';
    document.getElementById('ordersTable').innerHTML = html;
  } catch(e) {
    document.getElementById('ordersTable').innerHTML = '<div class="loading">Ошибка — сервер не отвечает на :8081</div>';
  }
}

function openStatusModal(orderID) {
  currentOrderID = orderID;
  document.getElementById('statusModal').classList.add('open');
}

async function confirmStatusChange() {
  const newStatus = document.getElementById('newStatus').value;
  closeModal('statusModal');
  try {
    const res = await fetch(API + '/admin/orders/' + currentOrderID + '/status', {
      method: 'PATCH', headers: {'Content-Type': 'application/json'},
      body: JSON.stringify({status: newStatus})
    });
    const data = await res.json();
    if (!res.ok) throw new Error(data.error);
    showToast('Статус обновлён: ' + (STATUS_RU[data.new_status] || data.new_status));
    loadOrders();
  } catch(e) { showToast(e.message, true); }
}

function openCourierModal(orderID) {
  currentOrderID = orderID;
  document.getElementById('newCourierID').innerHTML = '<option value="">— выберите курьера —</option>';
  document.getElementById('courierModal').classList.add('open');
  // Загружаем реальных свободных курьеров
  fetchT(COURIER_API + '/courier/couriers').then(r => r.json()).then(data => {
    const sel = document.getElementById('newCourierID');
    (data.couriers || []).forEach(c => {
      const opt = document.createElement('option');
      opt.value = c.id;
      opt.textContent = c.name + (c.status === 'busy' ? ' (занят)' : ' (свободен)');
      opt.disabled = c.status === 'busy';
      sel.appendChild(opt);
    });
  }).catch(() => {});
}

async function confirmCourierChange() {
  const courierID = document.getElementById('newCourierID').value;
  if (!courierID) { showToast('Выберите курьера из списка', true); return; }
  closeModal('courierModal');
  try {
    const res = await fetch(API + '/admin/orders/' + currentOrderID + '/courier', {
      method: 'PATCH', headers: {'Content-Type': 'application/json'},
      body: JSON.stringify({courier_id: courierID})
    });
    const data = await res.json();
    if (!res.ok) throw new Error(data.error);
    showToast('Курьер назначен');
    loadOrders();
  } catch(e) { showToast(e.message, true); }
}

async function cancelOrder(orderID) {
  if (!confirm('Отменить заказ?')) return;
  try {
    const res = await fetch(API + '/admin/orders/' + orderID + '/cancel', { method: 'POST' });
    const data = await res.json();
    if (!res.ok) throw new Error(data.error);
    showToast('Заказ отменён');
    loadOrders();
  } catch(e) { showToast(e.message, true); }
}

async function doRefund(orderID) {
  if (!confirm('Сделать возврат?')) return;
  try {
    const res = await fetch(API + '/admin/orders/' + orderID + '/refund', { method: 'POST' });
    const data = await res.json();
    if (!res.ok) throw new Error(data.error);
    showToast('Возврат выполнен');
    loadOrders();
  } catch(e) { showToast(e.message, true); }
}

// ─── ПОЛЬЗОВАТЕЛИ (с админского сервера) ─────────────────────────────────

async function loadUsers() {
  document.getElementById('usersTable').innerHTML = '<div class="loading">Загрузка...</div>';
  try {
    const res = await fetch(API + '/admin/users');
    if (!res.ok) throw new Error();
    const data = await res.json();
    const users = data.users || [];
    if (!users.length) {
      document.getElementById('usersTable').innerHTML = '<div class="loading">Пользователей пока нет — они появятся после оформления заказов</div>';
      return;
    }
    let html = '<table><thead><tr><th>ID</th><th>Имя</th><th>Email</th><th>Роль</th><th>Статус</th><th>Действия</th></tr></thead><tbody>';
    for (const u of users) {
      const sb = u.is_blocked
        ? '<span class="badge badge-blocked">Заблокирован</span>'
        : '<span class="badge badge-active">Активен</span>';
      const rb = '<span class="badge badge-' + (u.role || 'customer') + '">' + (u.role || 'customer') + '</span>';
      const ab = u.is_blocked
        ? '<button class="btn btn-success" onclick="toggleBlock(\'' + u.id + '\', false)">Разблокировать</button>'
        : '<button class="btn btn-danger" onclick="toggleBlock(\'' + u.id + '\', true)">Заблокировать</button>';
      html += '<tr>';
      html += '<td><code>' + u.id.slice(0,8) + '…</code></td>';
      html += '<td>' + (u.name || '—') + '</td>';
      html += '<td style="color:var(--muted)">' + (u.email || '—') + '</td>';
      html += '<td>' + rb + '</td>';
      html += '<td>' + sb + '</td>';
      html += '<td><div class="actions">' + ab + '</div></td>';
      html += '</tr>';
    }
    html += '</tbody></table>';
    document.getElementById('usersTable').innerHTML = html;
  } catch(e) {
    document.getElementById('usersTable').innerHTML = '<div class="loading">Ошибка — сервер не отвечает на :8081</div>';
  }
}

async function toggleBlock(userID, block) {
  const action = block ? 'block' : 'unblock';
  const label  = block ? 'заблокировать' : 'разблокировать';
  if (!confirm('Вы уверены, что хотите ' + label + ' пользователя?')) return;
  try {
    const res = await fetch(API + '/admin/users/' + userID + '/' + action, { method: 'PATCH' });
    const data = await res.json();
    if (!res.ok) throw new Error(data.error);
    showToast(block ? 'Пользователь заблокирован' : 'Пользователь разблокирован');
    loadUsers();
  } catch(e) { showToast(e.message, true); }
}
// ─── КУРЬЕРЫ (с курьерского сервера) ──────────────────────────────────────

async function loadCouriers() {
  document.getElementById('couriersTable').innerHTML = '<div class="loading">Загрузка...</div>';
  try {
    const res = await fetchT(COURIER_API + '/courier/couriers');
    if (!res.ok) throw new Error();
    const data = await res.json();
    const couriers = data.couriers || [];
    if (!couriers.length) {
      document.getElementById('couriersTable').innerHTML = '<div class="loading">Курьеров пока нет — они появятся после регистрации в приложении курьера</div>';
      return;
    }
    let html = '<table><thead><tr><th>ID</th><th>Имя</th><th>Телефон</th><th>Статус</th></tr></thead><tbody>';
    for (const c of couriers) {
      const badge = c.status === 'free'
        ? '<span class="badge badge-active">Свободен</span>'
        : '<span class="badge badge-assigned">Занят</span>';
      html += '<tr>';
      html += '<td><code>' + c.id.slice(0,8) + '…</code></td>';
      html += '<td>' + c.name + '</td>';
      html += '<td style="color:var(--muted)">' + (c.phone || '—') + '</td>';
      html += '<td>' + badge + '</td>';
      html += '</tr>';
    }
    html += '</tbody></table>';
    document.getElementById('couriersTable').innerHTML = html;
  } catch(e) {
    document.getElementById('couriersTable').innerHTML = '<div class="loading">Курьерский сервер не отвечает на :8082</div>';
  }
}

// ─── УТИЛИТЫ ──────────────────────────────────────────────────────────────

function closeModal(id) {
  document.getElementById(id).classList.remove('open');
}

function showToast(msg, isError = false) {
  const t = document.getElementById('toast');
  t.textContent = msg;
  t.className = 'toast show' + (isError ? ' error' : '');
  setTimeout(() => t.className = 'toast', 3000);
}

loadOrders();
</script>
</body>
</html>`
