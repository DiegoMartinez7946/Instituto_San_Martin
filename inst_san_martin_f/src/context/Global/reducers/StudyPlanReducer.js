export default (state, action) => {
  switch (action.type) {
    case 'ADD_STUDY_PLAN': {
      const prev = Array.isArray(state.studyPlans) ? state.studyPlans : [];
      return { ...state, studyPlans: [...prev, action.payload] };
    }
    case 'UPDATE_STUDY_PLAN': {
      const updated = (state.studyPlans || []).map((p) => {
        const pid = p.id || p.ID;
        const aid = action.payload.id || action.payload.ID;
        if (pid === aid) {
          return { ...p, ...action.payload };
        }
        return p;
      });
      return { ...state, studyPlans: updated };
    }
    case 'CHANGE_ACTIVE_STUDY_PLAN': {
      const changed = (state.studyPlans || []).map((p) => {
        const pid = p.id || p.ID;
        const aid = action.payload.id || action.payload.ID;
        if (pid === aid) {
          return { ...p, active: action.payload.active };
        }
        return p;
      });
      return { ...state, studyPlans: changed };
    }
    case 'GET_STUDY_PLAN':
      return { ...state, studyPlans: Array.isArray(action.payload) ? action.payload : [] };
    default:
      return state;
  }
};
