import { useMemo, useRef, useEffect } from 'react';
import './Heatmap.css';

const MONTHS = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];

export default function Heatmap({ days = [] }) {
  const scrollRef = useRef(null);

  useEffect(() => {
    if (scrollRef.current) {
      scrollRef.current.scrollLeft = scrollRef.current.scrollWidth;
    }
  }, [days]);

  const { weeks, monthLabels } = useMemo(() => {
    const map = {};
    for (const d of days) {
      const key = d.date.slice(0, 10);
      map[key] = (map[key] || 0) + d.quantity;
    }

    const values = Object.values(map);
    const max = Math.max(1, ...values);

    const today = new Date();
    const cells = [];
    const start = new Date(today);
    start.setDate(start.getDate() - 20 * 7);
    const dayOfWeek = start.getDay();
    const diff = dayOfWeek === 0 ? -6 : 1 - dayOfWeek;
    start.setDate(start.getDate() + diff);

    // Extend to end of current week (Saturday) so the grid is always full weeks
    const end = new Date(today);
    const endDow = end.getDay();
    end.setDate(end.getDate() + (endDow === 0 ? 0 : 7 - endDow));

    const cursor = new Date(start);
    while (cursor <= end) {
      const key = cursor.toISOString().slice(0, 10);
      cells.push({
        date: key,
        qty: map[key] || 0,
        level: map[key] ? Math.min(4, Math.ceil((map[key] / max) * 4)) : 0,
      });
      cursor.setDate(cursor.getDate() + 1);
    }

    const w = [];
    for (let i = 0; i < cells.length; i += 7) {
      w.push(cells.slice(i, i + 7));
    }

    const labels = [];
    let lastMonth = -1;
    for (let wi = 0; wi < w.length; wi++) {
      const week = w[wi];
      if (week.length === 0) continue;
      const month = new Date(week[0].date).getMonth();
      if (month !== lastMonth) {
        labels.push({ index: wi, label: MONTHS[month] });
        lastMonth = month;
      }
    }

    return { weeks: w, monthLabels: labels };
  }, [days]);

  const dayLabels = ['', 'Mon', '', 'Wed', '', 'Fri', ''];

  return (
    <div className="heatmap">
      <div className="heatmap-grid">
        <div className="heatmap-labels">
          <div className="heatmap-month-spacer" />
          {dayLabels.map((label, i) => (
            <div key={i} className="heatmap-day-label">{label}</div>
          ))}
        </div>
        <div className="heatmap-weeks" ref={scrollRef}>
          <div className="heatmap-month-row">
            {monthLabels.map((ml, i) => (
              <span
                key={i}
                className="heatmap-month-label"
                style={{ left: `calc(${ml.index} * var(--heatmap-col-w, 16px))` }}
              >
                {ml.label}
              </span>
            ))}
          </div>
          <div className="heatmap-cols">
            {weeks.map((week, wi) => (
              <div key={wi} className="heatmap-week">
                {week.map((cell, di) => (
                  <div
                    key={di}
                    className={`heatmap-cell${cell.level === 0 ? ' empty' : ''} level-${cell.level}`}
                    title={`${cell.date}: ${cell.qty}`}
                  />
                ))}
              </div>
            ))}
          </div>
        </div>
      </div>
      <div className="heatmap-legend">
        <span>Less</span>
        <div className="heatmap-cell level-0" />
        <div className="heatmap-cell level-1" />
        <div className="heatmap-cell level-2" />
        <div className="heatmap-cell level-3" />
        <div className="heatmap-cell level-4" />
        <span>More</span>
      </div>
    </div>
  );
}
