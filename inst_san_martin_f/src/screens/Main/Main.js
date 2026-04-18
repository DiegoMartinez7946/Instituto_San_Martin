import React, { useEffect, useState } from 'react';
import { Container } from 'react-bootstrap';

import clientAxios from '../../config/axios';
import { useGlobal } from '../../context/Global/GlobalProvider';
import { decodeTokenFromCookie, normalizeRoleFromToken } from '../../utils/jwt';

import styles from './Main.module.css';

const fmtInt = (n) =>
  n != null && Number.isFinite(Number(n))
    ? Number(n).toLocaleString('es-AR')
    : '—';

const fmtPct = (n) =>
  n != null && Number.isFinite(Number(n)) ? `${Number(n).toFixed(1)} %` : '—';

/** Load average (Linux/macOS): no es porcentaje. */
const fmtLoad = (n) =>
  n != null && Number.isFinite(Number(n)) ? Number(n).toFixed(2) : '—';

const fmtUptime = (sec) => {
  if (sec == null || !Number.isFinite(Number(sec))) return '—';
  const s = Number(sec);
  const d = Math.floor(s / 86400);
  const h = Math.floor((s % 86400) / 3600);
  const m = Math.floor((s % 3600) / 60);
  if (d > 0) return `${fmtInt(d)} d ${fmtInt(h)} h`;
  if (h > 0) return `${fmtInt(h)} h ${fmtInt(m)} min`;
  return `${fmtInt(m)} min`;
};

const fmtBytesMB = (n) => {
  if (n == null || !Number.isFinite(Number(n))) return '—';
  const mb = Number(n);
  return `${mb >= 100 ? mb.toFixed(0) : mb.toFixed(1)} MB`;
};

/** Tamaños devueltos por dbStats (bytes). */
const fmtBytes = (n) => {
  if (n == null || !Number.isFinite(Number(n))) return '—';
  const b = Number(n);
  if (b < 1024) return `${fmtInt(b)} B`;
  if (b < 1048576) return `${(b / 1024).toFixed(1)} KB`;
  if (b < 1073741824) return `${(b / 1048576).toFixed(2)} MB`;
  return `${(b / 1073741824).toFixed(2)} GB`;
};

const roleIsAdministrador = (globalState) => {
  const u = globalState.userLogin;
  if (u && typeof u === 'object' && Object.keys(u).length > 0) {
    return normalizeRoleFromToken(u) === 'ADMINISTRADOR';
  }
  return normalizeRoleFromToken(decodeTokenFromCookie()) === 'ADMINISTRADOR';
};

const MongoMonitorCard = ({ dbMetrics, dbMetricsErr }) => (
  <div className={styles.metricsCard}>
    <p className={styles.metricsCard__label}>MongoDB — monitor (solo administrador)</p>
    {dbMetricsErr ? (
      <div className={styles.metricsCard__warn}>{dbMetricsErr}</div>
    ) : null}
    {dbMetrics ? (
      <>
        {!dbMetrics.serverStatusOk && dbMetrics.serverStatusMsg ? (
          <div
            className={
              dbMetrics.partialMetricsOnly
                ? styles.metricsCard__info
                : styles.metricsCard__warn
            }
          >
            serverStatus: {dbMetrics.serverStatusMsg}
          </div>
        ) : null}
        <dl className={styles.metricsCard__grid}>
          <dt>Latencia ping</dt>
          <dd>{dbMetrics.pingMs != null ? `${dbMetrics.pingMs.toFixed(2)} ms` : '—'}</dd>
          <dt>Uptime servidor</dt>
          <dd>
            {dbMetrics.partialMetricsOnly
              ? '—'
              : dbMetrics.uptimeSec != null
              ? `${fmtInt(dbMetrics.uptimeSec)} s`
              : '—'}
          </dd>
        </dl>
        {dbMetrics.partialMetricsOnly && dbMetrics.dbStats ? (
          <div className={styles.metricsSection}>
            <p className={styles.metricsSectionTitle}>
              Base de datos (dbStats — métricas parciales)
            </p>
            <dl className={styles.metricsCard__grid}>
              <dt>Base</dt>
              <dd>{dbMetrics.dbStats.db || '—'}</dd>
              <dt>Colecciones</dt>
              <dd>{fmtInt(dbMetrics.dbStats.collections)}</dd>
              <dt>Documentos</dt>
              <dd>{fmtInt(dbMetrics.dbStats.objects)}</dd>
              <dt>Tamaño datos / almacenamiento</dt>
              <dd>
                {fmtBytes(dbMetrics.dbStats.dataSize)} /{' '}
                {fmtBytes(dbMetrics.dbStats.storageSize)}
              </dd>
              <dt>Índices (# / tamaño)</dt>
              <dd>
                {fmtInt(dbMetrics.dbStats.indexes)} /{' '}
                {fmtBytes(dbMetrics.dbStats.indexSize)}
              </dd>
            </dl>
          </div>
        ) : null}
        {!dbMetrics.partialMetricsOnly ? (
          <>
            <div className={styles.metricsSection}>
              <p className={styles.metricsSectionTitle}>Conexiones</p>
              <dl className={styles.metricsCard__grid}>
                <dt>Activas / disponibles</dt>
                <dd>
                  {fmtInt(dbMetrics.connections && dbMetrics.connections.current)} /{' '}
                  {fmtInt(dbMetrics.connections && dbMetrics.connections.available)}
                </dd>
                <dt>Creadas (total hist.)</dt>
                <dd>
                  {fmtInt(dbMetrics.connections && dbMetrics.connections.totalCreated)}
                </dd>
              </dl>
            </div>
            <div className={styles.metricsSection}>
              <p className={styles.metricsSectionTitle}>Memoria (proceso mongod)</p>
              <dl className={styles.metricsCard__grid}>
                <dt>Residente</dt>
                <dd>{fmtBytesMB(dbMetrics.memoryMb && dbMetrics.memoryMb.resident)}</dd>
                <dt>Virtual</dt>
                <dd>{fmtBytesMB(dbMetrics.memoryMb && dbMetrics.memoryMb.virtual)}</dd>
              </dl>
            </div>
            {(dbMetrics.wiredTigerCacheMb || 0) > 0 ? (
              <div className={styles.metricsSection}>
                <p className={styles.metricsSectionTitle}>WiredTiger caché</p>
                <dl className={styles.metricsCard__grid}>
                  <dt>Datos en caché (aprox.)</dt>
                  <dd>{fmtBytesMB(dbMetrics.wiredTigerCacheMb)}</dd>
                </dl>
              </div>
            ) : null}
            <div className={styles.metricsSection}>
              <p className={styles.metricsSectionTitle}>Operaciones (acumuladas)</p>
              <dl className={styles.metricsCard__grid}>
                <dt>Query</dt>
                <dd>{fmtInt(dbMetrics.opcounters && dbMetrics.opcounters.query)}</dd>
                <dt>Command</dt>
                <dd>{fmtInt(dbMetrics.opcounters && dbMetrics.opcounters.command)}</dd>
                <dt>Insert / Update / Delete</dt>
                <dd>
                  {fmtInt(dbMetrics.opcounters && dbMetrics.opcounters.insert)} /{' '}
                  {fmtInt(dbMetrics.opcounters && dbMetrics.opcounters.update)} /{' '}
                  {fmtInt(dbMetrics.opcounters && dbMetrics.opcounters.delete)}
                </dd>
              </dl>
            </div>
            <div className={styles.metricsSection}>
              <p className={styles.metricsSectionTitle}>Red (acumulado)</p>
              <dl className={styles.metricsCard__grid}>
                <dt>Bytes in / out</dt>
                <dd>
                  {fmtInt(dbMetrics.network && dbMetrics.network.bytesIn)} /{' '}
                  {fmtInt(dbMetrics.network && dbMetrics.network.bytesOut)}
                </dd>
                <dt>Peticiones #</dt>
                <dd>{fmtInt(dbMetrics.network && dbMetrics.network.numRequests)}</dd>
              </dl>
            </div>
          </>
        ) : null}
        {dbMetrics.note ? <p className={styles.metricsCard__note}>{dbMetrics.note}</p> : null}
      </>
    ) : !dbMetricsErr ? (
      <p className={styles.metricsCard__note}>Cargando métricas…</p>
    ) : null}
  </div>
);

const SystemMonitorCard = ({ systemMetrics, systemMetricsErr }) => (
  <div className={styles.metricsCard}>
    <p className={styles.metricsCard__label}>Servidor — API backend (solo administrador)</p>
    {systemMetricsErr ? (
      <div className={styles.metricsCard__warn}>{systemMetricsErr}</div>
    ) : null}
    {systemMetrics ? (
      <>
        <div className={styles.metricsSection}>
          <p className={styles.metricsSectionTitle}>Host</p>
          <dl className={styles.metricsCard__grid}>
            <dt>Nombre</dt>
            <dd>{systemMetrics.host && systemMetrics.host.hostname ? systemMetrics.host.hostname : '—'}</dd>
            <dt>Uptime SO</dt>
            <dd>{fmtUptime(systemMetrics.host && systemMetrics.host.uptimeSec)}</dd>
            <dt>Sistema</dt>
            <dd>
              {[systemMetrics.host && systemMetrics.host.os, systemMetrics.host && systemMetrics.host.platform]
                .filter(Boolean)
                .join(' · ') || '—'}
            </dd>
          </dl>
        </div>
        {systemMetrics.host &&
        systemMetrics.host.platform &&
        systemMetrics.host.platform !== 'windows' ? (
          <div className={styles.metricsSection}>
            <p className={styles.metricsSectionTitle}>Load average</p>
            <dl className={styles.metricsCard__grid}>
              <dt>1 / 5 / 15 min</dt>
              <dd>
                {fmtLoad(systemMetrics.loadAvg.load1)} ·{' '}
                {fmtLoad(systemMetrics.loadAvg.load5)} ·{' '}
                {fmtLoad(systemMetrics.loadAvg.load15)}
              </dd>
            </dl>
          </div>
        ) : null}
        <div className={styles.metricsSection}>
          <p className={styles.metricsSectionTitle}>Recursos</p>
          <dl className={styles.metricsCard__grid}>
            <dt>CPU (host)</dt>
            <dd>{fmtPct(systemMetrics.cpuPercent)}</dd>
            <dt>RAM usada</dt>
            <dd>{fmtPct(systemMetrics.memory && systemMetrics.memory.usedPercent)}</dd>
            <dt>Disco ({systemMetrics.disk && systemMetrics.disk.path ? systemMetrics.disk.path : '—'})</dt>
            <dd>{fmtPct(systemMetrics.disk && systemMetrics.disk.usedPercent)}</dd>
          </dl>
        </div>
        {systemMetrics.swap && systemMetrics.swap.totalBytes > 0 ? (
          <div className={styles.metricsSection}>
            <p className={styles.metricsSectionTitle}>Swap</p>
            <dl className={styles.metricsCard__grid}>
              <dt>Uso swap</dt>
              <dd>{fmtPct(systemMetrics.swap.usedPercent)}</dd>
            </dl>
          </div>
        ) : null}
        <div className={styles.metricsSection}>
          <p className={styles.metricsSectionTitle}>Proceso Go (esta API)</p>
          <dl className={styles.metricsCard__grid}>
            <dt>PID</dt>
            <dd>{fmtInt(systemMetrics.processSelf && systemMetrics.processSelf.pid)}</dd>
            <dt>Goroutines</dt>
            <dd>{fmtInt(systemMetrics.goRuntime && systemMetrics.goRuntime.goroutines)}</dd>
            <dt>Heap Go</dt>
            <dd>{fmtBytes(systemMetrics.goRuntime && systemMetrics.goRuntime.heapAllocBytes)}</dd>
            <dt>RSS proceso</dt>
            <dd>{fmtBytes(systemMetrics.processSelf && systemMetrics.processSelf.rssBytes)}</dd>
            <dt>CPU proceso</dt>
            <dd>{fmtPct(systemMetrics.processSelf && systemMetrics.processSelf.cpuPercent)}</dd>
          </dl>
        </div>
        {systemMetrics.topProcessesByRss && systemMetrics.topProcessesByRss.length > 0 ? (
          <div className={styles.metricsSection}>
            <p className={styles.metricsSectionTitle}>Top RAM (RSS en el host)</p>
            <table className={styles.processTable}>
              <thead>
                <tr>
                  <th>PID</th>
                  <th className={styles.cellName}>Proceso</th>
                  <th>RAM</th>
                  <th>CPU*</th>
                </tr>
              </thead>
              <tbody>
                {systemMetrics.topProcessesByRss.map((row) => (
                  <tr key={`${row.pid}-${row.name}`}>
                    <td>{fmtInt(row.pid)}</td>
                    <td className={styles.cellName} title={row.name}>
                      {row.name || '—'}
                    </td>
                    <td>{fmtBytes(row.rssBytes)}</td>
                    <td>{fmtPct(row.cpuPercent)}</td>
                  </tr>
                ))}
              </tbody>
            </table>
            <p className={styles.metricsCard__note}>
              * CPU por proceso puede mostrar 0 % en la primera muestra (instantáneo).
            </p>
          </div>
        ) : null}
        {systemMetrics.note ? <p className={styles.metricsCard__note}>{systemMetrics.note}</p> : null}
      </>
    ) : !systemMetricsErr ? (
      <p className={styles.metricsCard__note}>Cargando métricas del servidor…</p>
    ) : null}
  </div>
);

const Main = () => {
  const [globalState] = useGlobal();
  const [dbMetrics, setDbMetrics] = useState(null);
  const [dbMetricsErr, setDbMetricsErr] = useState(null);
  const [systemMetrics, setSystemMetrics] = useState(null);
  const [systemMetricsErr, setSystemMetricsErr] = useState(null);

  const t = globalState.portalTime;
  const showDbMonitor = roleIsAdministrador(globalState);

  useEffect(() => {
    if (!showDbMonitor) {
      setDbMetrics(null);
      setDbMetricsErr(null);
      setSystemMetrics(null);
      setSystemMetricsErr(null);
      return undefined;
    }

    const authHeader = () => {
      const token = typeof document !== 'undefined' ? document.cookie.replace('token=', '') : '';
      const trimmed = token.trim();
      return trimmed ? { Authorization: `Bearer${trimmed}` } : {};
    };

    const fetchMetrics = async () => {
      const headers = authHeader();
      if (!headers.Authorization) return;

      const runDb = async () => {
        try {
          const res = await clientAxios.get('/admin/db-metrics', { headers });
          const payload = res.data && res.data.data;
          const code = res.data && res.data.code !== undefined ? Number(res.data.code) : 200;
          if (code >= 400) {
            setDbMetricsErr(res.data.message || `Error ${code}`);
            setDbMetrics(null);
            return;
          }
          setDbMetrics(payload || null);
          setDbMetricsErr(null);
        } catch (e) {
          let errMsg = 'No se pudieron obtener las métricas de MongoDB';
          if (e && e.message) errMsg = e.message;
          if (e && e.response && e.response.data && e.response.data.message) {
            errMsg = e.response.data.message;
          }
          setDbMetricsErr(errMsg);
          setDbMetrics(null);
        }
      };

      const runSystem = async () => {
        try {
          const res = await clientAxios.get('/admin/system-metrics', { headers });
          const payload = res.data && res.data.data;
          const code = res.data && res.data.code !== undefined ? Number(res.data.code) : 200;
          if (code >= 400) {
            setSystemMetricsErr(res.data.message || `Error ${code}`);
            setSystemMetrics(null);
            return;
          }
          setSystemMetrics(payload || null);
          setSystemMetricsErr(null);
        } catch (e) {
          let errMsg = 'No se pudieron obtener las métricas del servidor';
          if (e && e.message) errMsg = e.message;
          if (e && e.response && e.response.data && e.response.data.message) {
            errMsg = e.response.data.message;
          }
          setSystemMetricsErr(errMsg);
          setSystemMetrics(null);
        }
      };

      await Promise.all([runDb(), runSystem()]);
    };

    fetchMetrics();
    const id = setInterval(fetchMetrics, 50000);
    return () => clearInterval(id);
  }, [showDbMonitor]);

  return (
    <Container className={styles.container}>
      <div className={styles.titleRow}>
        <h1 className={styles.title}>home</h1>
      </div>

      {showDbMonitor ? (
        <div className={styles.dashboardOuter}>
          <div className={styles.dashboardRow}>
            <div className={styles.dashboardColLeft}>
              <div className={styles.dashboardLeftStack}>
                {t ? (
                  <div className={styles.clockCard}>
                    <p className={styles.clockCard__label}>Hora en Argentina</p>
                    <p className={styles.clockCard__hora}>{t.hora}</p>
                    <dl className={styles.clockCard__grid}>
                      <dt>Fecha</dt>
                      <dd>{t.fecha}</dd>
                      <dt>Día</dt>
                      <dd>{t.dia}</dd>
                      <dt>Mes</dt>
                      <dd>{t.mes}</dd>
                      <dt>Año</dt>
                      <dd>{t.año}</dd>
                    </dl>
                  </div>
                ) : null}
                <div className={styles.dashboardMongoGrow}>
                  <MongoMonitorCard dbMetrics={dbMetrics} dbMetricsErr={dbMetricsErr} />
                </div>
              </div>
            </div>
            <div className={styles.dashboardColRight}>
              <SystemMonitorCard systemMetrics={systemMetrics} systemMetricsErr={systemMetricsErr} />
            </div>
          </div>
        </div>
      ) : t ? (
        <div className={styles.dashboardOuter}>
          <div className={`${styles.dashboardRow} ${styles.dashboardRowOneCol}`}>
            <div className={styles.dashboardColLeft}>
              <div className={styles.clockCard}>
                <p className={styles.clockCard__label}>Hora en Argentina</p>
                <p className={styles.clockCard__hora}>{t.hora}</p>
                <dl className={styles.clockCard__grid}>
                  <dt>Fecha</dt>
                  <dd>{t.fecha}</dd>
                  <dt>Día</dt>
                  <dd>{t.dia}</dd>
                  <dt>Mes</dt>
                  <dd>{t.mes}</dd>
                  <dt>Año</dt>
                  <dd>{t.año}</dd>
                </dl>
              </div>
            </div>
          </div>
        </div>
      ) : null}
    </Container>
  );
};

export default Main;
