/**
 * Utilitas bersama untuk tipe Job Value — single source of truth label tipe.
 * Dipakai oleh: AppLayout, HeaderBar, JobValuesIndex, JobValueSection, JobValuesForm.
 *
 * Ketika menambah tipe job value baru, update dua hal di sini:
 *   1. VALID_JOB_VALUE_TYPES
 *   2. JOB_VALUE_TYPE_LABEL_FALLBACK
 * plus key locale `job_values.types.*` dan `job_values.type_desc.*`.
 */

// Daftar tipe job value yang valid (urut sesuai tampilan index card)
export const VALID_JOB_VALUE_TYPES = [
  'education', 'experience', 'subordinate', 'activity', 'communicating_influencing_skill',
  'thinking_environment', 'thinking_chalenge', 'kecerdasan', 'innovation_creativity',
  'self_confidence', 'flexibility', 'tenacity', 'continuous_learning',
  'competency_based_human_resources_management', 'competency_development', 'people_development',
  'career_management', 'hr_assessment', 'recruitement_selection', 'job_analysis_evaluation',
  'organizational_development', 'human_resources_information_system', 'workload_analysis',
  'performance_apraisal', 'remuneration_manajemen', 'reward_punisment_management',
  'health_safety_environment', 'hubungan_industrial', 'budgeting',
  'integrity', 'achievement_orientation', 'building_partnership', 'planning_organizing',
  'leadership', 'developing_others',
  'environment', 'risk', 'relationship', 'frequency', 'asset', 'asset_authority',
  'authority', 'authority_unauthorized', 'cash', 'impact', 'impact_unauthorized'
]

// Fallback label bila key locale belum ada (t() bisa return key string saat missing)
export const JOB_VALUE_TYPE_LABEL_FALLBACK = {
  education: 'Education',
  experience: 'Experience',
  activity: 'Physical Activity',
  communicating_influencing_skill: 'Communicating & Influencing Skill',
  thinking_environment: 'Thinking Environment',
  thinking_chalenge: 'Thinking Challenge',
  kecerdasan: 'Intelligence',
  innovation_creativity: 'Innovation & Creativity',
  self_confidence: 'Self Confidence',
  flexibility: 'Flexibility',
  tenacity: 'Tenacity',
  continuous_learning: 'Continuous Learning',
  competency_based_human_resources_management: 'Competency Based HR Management',
  competency_development: 'Competency Development',
  people_development: 'People Development',
  career_management: 'Career Management',
  hr_assessment: 'HR Assessment',
  recruitement_selection: 'Recruitment & Selection',
  job_analysis_evaluation: 'Job Analysis & Evaluation',
  organizational_development: 'Organizational Development',
  human_resources_information_system: 'HR Information System',
  workload_analysis: 'Workload Analysis',
  performance_apraisal: 'Performance Appraisal',
  remuneration_manajemen: 'Remuneration Management',
  reward_punisment_management: 'Reward & Punishment Management',
  health_safety_environment: 'Health & Safety Environment',
  hubungan_industrial: 'Industrial Relations',
  budgeting: 'Budgeting',
  integrity: 'Integrity',
  achievement_orientation: 'Achievement Orientation',
  building_partnership: 'Building Partnership',
  planning_organizing: 'Planning & Organizing',
  leadership: 'Leadership',
  developing_others: 'Developing Others',
  environment: 'Environment',
  risk: 'Risk',
  relationship: 'Relationship',
  frequency: 'Frequency',
  asset: 'Asset',
  asset_authority: 'Asset Authority',
  authority: 'Authority',
  authority_unauthorized: 'Authority (No Financial Authority)',
  cash: 'Cash',
  impact: 'Impact',
  impact_unauthorized: 'Impact (No Financial Authority)',
  subordinate: 'Total Subordinates'
}

/** Ambil label tipe (bilingual, fallback ke map lokal bila key locale missing). */
export function jobValueTypeLabel(t, type) {
  if (!type) return ''
  const label = t(`job_values.types.${type}`)
  if (label && !label.startsWith('job_values.types.')) return label
  return JOB_VALUE_TYPE_LABEL_FALLBACK[type] || type
}

/** Ambil deskripsi tipe (bilingual, fallback '' bila key locale missing). */
export function jobValueTypeDesc(t, type) {
  if (!type) return ''
  const key = `job_values.type_desc.${type}`
  const desc = t(key)
  if (desc !== key) return desc
  return ''
}

/** Normalisasi :type dari URL — kembalikan tipe valid atau fallback. */
export function normalizeJobValueType(raw) {
  return VALID_JOB_VALUE_TYPES.includes(raw) ? raw : 'education'
}
