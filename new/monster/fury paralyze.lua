	local combat = Combat()
	combat:setParameter(COMBAT_PARAM_EFFECT, CONST_ME_MORTAREA)
	combat:setParameter(COMBAT_PARAM_DISTANCEEFFECT, CONST_ANI_THROWINGSTAR)

	local condition = Condition(CONDITION_PARALYZE)
	condition:setParameter(CONDITION_PARAM_TICKS, 10000)
	condition:setFormula(-0.25, 0, -0.35, 0)
	combat:addCondition(condition)

local spell = Spell("instant")

function spell.onCastSpell(creature, var)
	return combat:execute(creature, var)
end

spell:group("support")
spell:spellid("194")
spell:name("fury paralyze")
spell:words("###168")
spell:lvl("50")
spell:mana("30")
spell:range("3")
spell:cooldown("40000")
spell:groupcooldown("2000")
spell:isPremium(true)
spell:isAggressive(true)
spell:blockWalls(true)
spell:needTarget(true)
spell:needLearn(true)
spell:needDirection(true)
spell:isSelfTarget(true)
spell:register()