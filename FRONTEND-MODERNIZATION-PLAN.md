# Frontend Test Infrastructure Modernization Plan

## Current State Analysis (December 2024)

### Package Versions - Outdated but Functional

**Angular Ecosystem** (Last updated: 2021)
- Angular: v12.0.2 (Current: v21.0.5, Released 3.5 years ago)
- TypeScript: 4.2.4 (Current: 5.7.x)
- RxJS: 6.5.3 (Current: 7.8.2)

**Electron** (Security Critical)
- Current: 15.2.0 (Released Oct 2021)
- Latest in electron/package.json: 15.2.0
- Dependabot tried to bump to 22.3.25 (not merged)
- Latest Stable: 33.x
- Latest LTS: 28.x

**Testing Tools**
- Karma: 6.3.2 (still maintained)
- Jasmine: 3.10.0 (recently updated by us, was 3.6.0)
- Protractor: 7.0.0 ⚠️ **DEPRECATED** - Angular removed official support

**Node/NPM**
- CI uses: Node 20 ✅ (Latest LTS)
- npm ci ✅ (Modern practice)

### CI Configuration - Modern

**GitHub Actions** (.github/workflows/tests.yml)
- ✅ Uses actions/setup-go@v5 (latest)
- ✅ Uses actions/setup-node@v4 (latest)
- ✅ Uses actions/checkout@v4 (latest)
- ✅ Node 20 (current LTS)
- ✅ Go 1.24.x (current)
- ✅ Proper caching strategy

**Test Matrix**
- units
- integrations
- integrations/disable-csrf
- integrations/auth

---

## Recommended Updates

### Priority 1: Security Critical (Do Now)

#### 1. Electron Security Update
**Issue**: Electron 15.2.0 has known security vulnerabilities

**Options**:
- **Conservative**: Electron 22.3.27 (same as Dependabot attempted)
  - Minimal breaking changes from v15
  - Security patches included
  - Still in extended support
  
- **Recommended**: Electron 28.x (Current LTS)
  - Long-term support
  - Better security
  - May require minor API updates

**Action**: Update electron/package.json
```json
{
  "devDependencies": {
    "electron": "^28.3.3",
    "electron-builder": "^25.1.8"
  }
}
```

**Testing Required**: Build wallet and verify:
- App launches
- IPC communication works
- Window management works
- File dialogs work

#### 2. Remove Protractor (Deprecated)
**Issue**: Protractor is officially deprecated by Angular team

**Recommended Replacement**: 
- **Playwright** (recommended by Angular team)
- **Cypress** (popular alternative)

**Action**: 
```json
// Remove from package.json
"protractor": "~7.0.0",  // DELETE

// Add (choose one):
"@playwright/test": "^1.48.0",  // Recommended
// OR
"cypress": "^13.16.0",
```

**Migration Effort**: ~2-4 weeks to rewrite E2E tests

**Alternative**: Keep protractor for now, skip in CI with comment explaining deprecated status

---

### Priority 2: Minor Security Patches (Safe to Update)

#### Package Updates (Same Major Version)
```json
// src/gui/static/package.json updates
{
  "dependencies": {
    "moment": "^2.30.1",  // Currently: 2.21.0 (security fixes)
    "core-js": "^2.6.12", // Latest v2.x (v3 would break Angular 12)
    "bootstrap": "^4.6.2", // Latest v4.x (v5 would break)
  },
  "devDependencies": {
    "karma": "~6.4.4",    // Currently: 6.3.2
    "typescript": "~4.2.4", // Keep - Angular 12 compatible
  }
}
```

**Risk**: Very low - same major versions
**Testing**: Run existing test suite

---

### Priority 3: Test Infrastructure Improvements

#### 1. Add npm-check-updates to CI
```yaml
# .github/workflows/tests.yml - add new job
  Dependency-Audit:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: 20
      - name: Audit dependencies
        run: |
          cd src/gui/static
          npm audit --audit-level=moderate
```

#### 2. Pin Node version more strictly
```yaml
# Current: node-version: 20 (any 20.x)
# Better:  node-version: 20.18.1 (exact LTS)
```

#### 3. Add package-lock.json validation
```yaml
- name: Verify package-lock.json
  run: |
    cd src/gui/static
    npm ci
    git diff --exit-code package-lock.json
```

#### 4. Cache npm dependencies
```yaml
- name: Cache node modules
  uses: actions/cache@v4
  with:
    path: ~/.npm
    key: ${{ runner.os }}-node-${{ hashFiles('**/package-lock.json') }}
```

---

### Priority 4: Long-term (Future Work)

#### 1. Angular Upgrade Path
**Current**: Angular 12 (EOL: November 2022)
**Target**: Angular 17+ (current LTS)

**Migration Steps** (each is a major effort):
- Angular 12 → 13 (TypeScript 4.4+)
- Angular 13 → 14 (TypeScript 4.6+)
- Angular 14 → 15 (TypeScript 4.8+, standalone components)
- Angular 15 → 16 (TypeScript 4.9+)
- Angular 16 → 17 (TypeScript 5.2+, new control flow)

**Estimated Effort**: 3-6 months
**Breaking Changes**: Major
**Benefit**: Security, performance, modern features

**Recommendation**: Plan separately, not in this PR

#### 2. Replace Protractor with Playwright
**Effort**: 2-4 weeks
**Benefits**: 
- Actively maintained
- Faster execution
- Better debugging
- Cross-browser support

**Recommendation**: Plan separately after Angular upgrade

---

## Immediate Action Plan (This PR)

### Phase 1: Safe Security Updates ✅ Can Do Now

```bash
# 1. Update electron (electron/package.json)
npm install --save-dev electron@^28.3.3 electron-builder@^25.1.8

# 2. Update safe dependencies (src/gui/static/package.json)
cd src/gui/static
npm install moment@^2.30.1
npm install --save-dev karma@~6.4.4

# 3. Regenerate lockfiles
npm install
cd ../../electron
npm install

# 4. Test locally
cd ../..
make install-deps-ui
make test-ui
make build-ui-travis
```

### Phase 2: CI Improvements ✅ Can Do Now

1. Add dependency audit job to workflow
2. Add npm cache to speed up CI
3. Pin Node to exact LTS version (20.18.1)

### Phase 3: Documentation ⚠️ Do Now

Add to README or CONTRIBUTING.md:
```markdown
## Frontend Dependencies

### Current Versions (Dec 2024)
- Angular: 12.0.2 (⚠️ EOL - upgrade planned)
- Electron: 28.3.3 ✅
- Node: 20.18.1 (LTS) ✅
- Protractor: DEPRECATED - migration to Playwright planned

### Upgrade Policy
- Security patches: Apply immediately
- Minor versions: Review and test
- Major versions: Requires RFC and planning

### Known Technical Debt
1. Angular 12 → 17+ migration (planned Q1 2025)
2. Protractor → Playwright migration (planned Q2 2025)
3. Bootstrap 4 → 5 migration (blocked on Angular upgrade)
```

---

## What NOT to Do (Right Now)

❌ **Don't upgrade Angular beyond v12** without a dedicated migration effort
❌ **Don't remove Protractor** until replacement E2E tests are written
❌ **Don't upgrade to Node 22** (not LTS yet, stick with 20.x)
❌ **Don't upgrade TypeScript beyond 4.2.x** (Angular 12 incompatible with 4.3+)
❌ **Don't upgrade Bootstrap to v5** (breaking changes, needs Angular upgrade first)
❌ **Don't upgrade RxJS to v7** without Angular upgrade (compatibility issues)

---

## Testing Checklist

After any updates:

### UI Tests
- [ ] `make install-deps-ui` - No errors
- [ ] `make lint-ui` - Passes
- [ ] `make test-ui` - All tests pass
- [ ] `make test-ui-e2e` - E2E tests pass
- [ ] `make build-ui-travis` - Build succeeds

### Integration Tests
- [ ] `make integration-test-stable`
- [ ] `make integration-test-stable-disable-csrf`
- [ ] `make integration-test-stable-auth`

### Electron Build
- [ ] Wallet app launches
- [ ] No console errors
- [ ] Basic operations work (create wallet, send, receive)

### Cross-platform
- [ ] Test on Linux
- [ ] Test on macOS (if available)
- [ ] Test on Windows (if available)

---

## Risk Assessment

| Update | Risk | Impact | Effort |
|--------|------|--------|--------|
| Electron 15→28 | Medium | High | 1-2 days |
| electron-builder | Low | Medium | 1 day |
| moment, karma patches | Very Low | Low | 1 hour |
| CI improvements | Very Low | Low | 2 hours |
| Remove Protractor | High | High | 2-4 weeks |
| Angular 12→17 | Very High | Very High | 3-6 months |

**Recommendation**: 
1. Do Electron + minor patches now (this PR)
2. CI improvements now (this PR)
3. Plan Angular upgrade separately (Q1 2025)
4. Plan Protractor replacement after Angular upgrade (Q2 2025)

---

## Files to Update (This PR)

1. `electron/package.json` - Electron 28.3.3, electron-builder 25.1.8
2. `electron/package-lock.json` - Regenerated
3. `src/gui/static/package.json` - moment, karma patches
4. `src/gui/static/package-lock.json` - Regenerated
5. `.github/workflows/tests.yml` - Add audit, caching, pin Node version
6. `CONTRIBUTING.md` or `README.md` - Document current state and upgrade plans

---

## Success Criteria

✅ All CI tests pass
✅ No new security vulnerabilities (npm audit clean)
✅ Electron wallet builds on all platforms
✅ No regression in functionality
✅ CI runs faster (with caching)
✅ Documentation updated

---

## References

- [Angular Update Guide](https://update.angular.io/)
- [Electron Releases](https://www.electronjs.org/docs/latest/tutorial/electron-timelines)
- [Protractor Deprecation Notice](https://github.com/angular/protractor/issues/5502)
- [GitHub Actions Best Practices](https://docs.github.com/en/actions/learn-github-actions/workflow-syntax-for-github-actions)
