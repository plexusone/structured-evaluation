package html

import "html/template"

// claimsTemplate is the self-contained HTML document for a claims report.
// All CSS is inline so the output has no external dependencies.
var claimsTemplate = template.Must(template.New("claims").Parse(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>{{.Title}} — Claims Audit</title>
<style>
  :root {
    --bg: #ffffff; --fg: #1a1a1a; --muted: #666; --border: #e2e2e2;
    --verified: #047857; --needs-review: #b45309; --rejected: #b91c1c; --unverified: #6b7280;
    --verified-bg: #ecfdf5; --needs-review-bg: #fffbeb; --rejected-bg: #fef2f2; --unverified-bg: #f9fafb;
  }
  @media (prefers-color-scheme: dark) {
    :root {
      --bg: #16181c; --fg: #e8e8e8; --muted: #9aa0a6; --border: #2c2f36;
      --verified: #34d399; --needs-review: #fbbf24; --rejected: #f87171; --unverified: #9ca3af;
      --verified-bg: #0c2a20; --needs-review-bg: #2a2110; --rejected-bg: #2a1414; --unverified-bg: #1f2126;
    }
  }
  * { box-sizing: border-box; }
  body { margin: 0; background: var(--bg); color: var(--fg);
    font: 15px/1.55 -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif; }
  .wrap { max-width: 920px; margin: 0 auto; padding: 2rem 1.25rem 4rem; }
  header h1 { font-size: 1.5rem; margin: 0 0 .25rem; }
  .meta { color: var(--muted); font-size: .85rem; margin-bottom: 1.25rem; }
  .meta code { font-size: .8rem; }
  .decision { display: inline-flex; align-items: center; gap: .5rem; font-weight: 700;
    padding: .4rem .8rem; border-radius: 8px; margin-bottom: 1rem; }
  .decision.pass { background: var(--verified-bg); color: var(--verified); }
  .decision.conditional { background: var(--needs-review-bg); color: var(--needs-review); }
  .decision.fail { background: var(--rejected-bg); color: var(--rejected); }
  .decision small { font-weight: 400; opacity: .85; }
  .chips { display: flex; flex-wrap: wrap; gap: .5rem; margin: 0 0 1.75rem; }
  .chip { border: 1px solid var(--border); border-radius: 999px; padding: .3rem .7rem; font-size: .82rem; }
  .chip b { font-variant-numeric: tabular-nums; }
  .chip.verified { color: var(--verified); }
  .chip.needs-review { color: var(--needs-review); }
  .chip.rejected { color: var(--rejected); }
  .chip.unverified { color: var(--unverified); }
  section.group { margin-bottom: 2rem; }
  section.group > h2 { font-size: 1rem; text-transform: uppercase; letter-spacing: .04em;
    display: flex; align-items: center; gap: .5rem; padding-bottom: .4rem; border-bottom: 2px solid var(--border); }
  .group.verified > h2 { color: var(--verified); }
  .group.needs-review > h2 { color: var(--needs-review); }
  .group.rejected > h2 { color: var(--rejected); }
  .group.unverified > h2 { color: var(--unverified); }
  .claim { border: 1px solid var(--border); border-left: 4px solid var(--border);
    border-radius: 8px; padding: .9rem 1rem; margin: .75rem 0; background: var(--bg); }
  .claim.verified { border-left-color: var(--verified); background: var(--verified-bg); }
  .claim.needs-review { border-left-color: var(--needs-review); background: var(--needs-review-bg); }
  .claim.rejected { border-left-color: var(--rejected); background: var(--rejected-bg); }
  .claim.unverified { border-left-color: var(--unverified); background: var(--unverified-bg); }
  .claim-head { display: flex; align-items: baseline; gap: .5rem; }
  .claim-icon { font-weight: 700; }
  .claim.verified .claim-icon { color: var(--verified); }
  .claim.needs-review .claim-icon { color: var(--needs-review); }
  .claim.rejected .claim-icon { color: var(--rejected); }
  .claim.unverified .claim-icon { color: var(--unverified); }
  .claim-text { font-weight: 600; }
  .claim-id { margin-left: auto; color: var(--muted); font-size: .72rem;
    font-family: ui-monospace, SFMono-Regular, Menlo, monospace; }
  .stat { margin: .5rem 0 .25rem; font-size: .9rem; }
  .stat .val { font-weight: 700; font-variant-numeric: tabular-nums; }
  .tag { display: inline-block; font-size: .72rem; padding: .1rem .45rem; border-radius: 4px;
    border: 1px solid var(--border); margin-left: .35rem; color: var(--muted); }
  .rel-high { color: var(--verified); border-color: var(--verified); }
  .rel-medium { color: var(--needs-review); border-color: var(--needs-review); }
  .rel-low { color: var(--rejected); border-color: var(--rejected); }
  .src { font-size: .82rem; margin: .35rem 0; }
  .src a { color: inherit; }
  .quote { font-size: .82rem; color: var(--muted); border-left: 2px solid var(--border);
    padding-left: .6rem; margin: .4rem 0; font-style: italic; }
  .rationale { font-size: .88rem; margin-top: .5rem; }
  .rationale b { color: var(--muted); font-weight: 600; text-transform: uppercase; font-size: .7rem;
    letter-spacing: .04em; display: block; margin-bottom: .15rem; }
  footer { margin-top: 3rem; color: var(--muted); font-size: .78rem; border-top: 1px solid var(--border); padding-top: 1rem; }
  @media print { body { background: #fff; color: #000; } .claim { break-inside: avoid; } }
</style>
</head>
<body>
<div class="wrap">
<header>
  <h1>{{.Title}}</h1>
  <div class="meta">
    Claims audit{{if .GeneratedAt}} · generated {{.GeneratedAt}}{{end}}{{if .ValidatedBy}} · {{.ValidatedBy}}{{end}}
    {{if .Document}}<br><code>{{.Document}}</code>{{end}}
  </div>
  <div class="decision {{.DecisionClass}}">{{.DecisionLabel}}{{if .DecisionText}} <small>— {{.DecisionText}}</small>{{end}}</div>
  <div class="chips">
    <span class="chip"><b>{{.Counts.Total}}</b> total</span>
    <span class="chip verified"><b>{{.Counts.Verified}}</b> verified</span>
    <span class="chip needs-review"><b>{{.Counts.NeedsReview}}</b> needs review</span>
    <span class="chip rejected"><b>{{.Counts.Rejected}}</b> rejected</span>
    {{if .Counts.Unverified}}<span class="chip unverified"><b>{{.Counts.Unverified}}</b> unverified</span>{{end}}
  </div>
</header>

{{range .Groups}}
<section class="group {{.Class}}">
  <h2><span>{{.Icon}}</span> {{.Label}} <span class="claim-id">{{len .Claims}}</span></h2>
  {{range .Claims}}
  <article class="claim {{.Class}}">
    <div class="claim-head">
      <span class="claim-icon">{{.Icon}}</span>
      <span class="claim-text">{{.Text}}</span>
      <span class="claim-id">{{.ID}}</span>
    </div>
    {{if .HasStat}}
    <div class="stat">
      <span class="val">{{.Value}}{{if .Unit}} {{.Unit}}{{end}}</span>
      {{if .Precision}}<span class="tag">{{.Precision}}</span>{{end}}
      {{if .AsOfDate}}<span class="tag">as of {{.AsOfDate}}</span>{{end}}
    </div>
    {{end}}
    {{if .HasSource}}
    <div class="src">
      Source: {{if .URL}}<a href="{{.URL}}" rel="noopener noreferrer">{{.URL}}</a>{{else}}<em>none</em>{{end}}
      {{if .SourceType}}<span class="tag">{{.SourceType}}</span>{{end}}
      {{if .Reliability}}<span class="tag {{.RelClass}}">{{.Reliability}}</span>{{end}}
    </div>
    {{if .QuotedText}}<div class="quote">“{{.QuotedText}}”</div>{{end}}
    {{end}}
    {{if .IsDerived}}<div class="src">Derived from: <code>{{.DerivedFrom}}</code></div>{{end}}
    {{if .Rationale}}<div class="rationale"><b>Commentary</b>{{.Rationale}}</div>{{end}}
  </article>
  {{end}}
</section>
{{end}}

<footer>
  Generated by structured-evaluation (render/html). Verdicts: <b>verified</b> = publishable;
  <b>needs review</b> = present with caveats / human decision; <b>rejected</b> = do not use.
</footer>
</div>
</body>
</html>
`))
